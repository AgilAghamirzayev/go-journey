package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

//
// ===== 1) Core CMS primitives: Base + Content Types =====
//

// Base is embedded into every content type.
// Reflection will set ID/Slug/Timestamps dynamically.
type Base struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Example content type #1
type Article struct {
	Base
	Title string   `json:"title" cms:"slug" required:"true"` // cms:"slug" -> source for slug
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

// Example content type #2
type Event struct {
	Base
	Name      string    `json:"name" cms:"slug" required:"true"`
	StartsAt  time.Time `json:"startsAt"`
	Location  string    `json:"location"`
	IsOnline  bool      `json:"isOnline"`
	MaxSeats  int       `json:"maxSeats"`
	EventCode string    `json:"eventCode"`
}

//
// ===== 2) Registry: reflect.Type per content name =====
//

type typeReg struct {
	Name string
	Type reflect.Type
}

var registry = map[string]typeReg{}
var registryMu sync.RWMutex

// RegisterContent[T] registers a struct type for dynamic routing + storage.
func RegisterContent[T any](name string) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("RegisterContent only supports struct types")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = typeReg{Name: name, Type: t}
}

//
// ===== 3) In-memory storage (replaceable by any backend later) =====
//

type Store struct {
	mu    sync.RWMutex
	data  map[string]map[string]any    // type -> id -> instance (pointer to struct)
	slugs map[string]map[string]string // type -> slug -> id
	seq   uint64
}

func NewStore() *Store {
	return &Store{
		data:  make(map[string]map[string]any),
		slugs: make(map[string]map[string]string),
	}
}

func (s *Store) newID() string {
	s.seq++
	return fmt.Sprintf("%d%03d", time.Now().UnixNano(), s.seq)
}

func (s *Store) slugExists(typ, slug string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.slugs[typ] == nil {
		return false
	}
	_, ok := s.slugs[typ][slug]
	return ok
}

func (s *Store) Create(typ string, objPtr any) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[typ] == nil {
		s.data[typ] = map[string]any{}
	}
	if s.slugs[typ] == nil {
		s.slugs[typ] = map[string]string{}
	}

	v := reflect.Indirect(reflect.ValueOf(objPtr))
	id := getStringField(v, "ID")
	if id == "" {
		id = s.newID()
		setField(v, "ID", id)
	}
	// maintain slug index
	if slug := getStringField(v, "Slug"); slug != "" {
		if _, exists := s.slugs[typ][slug]; exists {
			return nil, fmt.Errorf("slug already exists")
		}
		s.slugs[typ][slug] = id
	}

	s.data[typ][id] = objPtr
	return objPtr, nil
}

func (s *Store) Get(typ, id string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	obj, ok := s.data[typ][id]
	return obj, ok
}

func (s *Store) GetBySlug(typ, slug string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.slugs[typ][slug]
	if !ok {
		return nil, false
	}
	obj, ok := s.data[typ][id]
	return obj, ok
}

func (s *Store) List(typ string) []any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]any, 0, len(s.data[typ]))
	for _, v := range s.data[typ] {
		out = append(out, v)
	}
	return out
}

func (s *Store) Delete(typ, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, ok := s.data[typ][id]
	if !ok {
		return false
	}
	v := reflect.Indirect(reflect.ValueOf(obj))
	if slug := getStringField(v, "Slug"); slug != "" {
		delete(s.slugs[typ], slug)
	}
	delete(s.data[typ], id)
	return true
}

func (s *Store) Update(typ, id string, patch map[string]any, reg typeReg) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, ok := s.data[typ][id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	v := reflect.Indirect(reflect.ValueOf(obj))
	t := reg.Type

	// Track slug changes to maintain index
	var oldSlug string
	if slug := getStringField(v, "Slug"); slug != "" {
		oldSlug = slug
	}

	for key, raw := range patch {
		idx, field, found := findFieldByJSON(t, key)
		if !found {
			// Allow patch by Go field name too
			idx2, field2, found2 := findFieldByName(t, key)
			if !found2 {
				continue
			}
			idx, field = idx2, field2
		}
		fv := v.FieldByIndex(idx)
		cv, err := convertJSONValue(raw, fv.Type())
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}
		fv.Set(cv)
	}

	setField(v, "UpdatedAt", time.Now().UTC())

	// update slug index if changed
	newSlug := getStringField(v, "Slug")
	if newSlug != oldSlug {
		// remove old
		if oldSlug != "" && s.slugs[typ][oldSlug] == id {
			delete(s.slugs[typ], oldSlug)
		}
		// add new
		if newSlug != "" {
			if _, exists := s.slugs[typ][newSlug]; exists {
				return nil, fmt.Errorf("slug already exists")
			}
			if s.slugs[typ] == nil {
				s.slugs[typ] = map[string]string{}
			}
			s.slugs[typ][newSlug] = id
		}
	}

	return obj, nil
}

//
// ===== 4) Reflection helpers =====
//

func setField(v reflect.Value, field string, value any) {
	f := v.FieldByName(field)
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(value))
	}
}

func getStringField(v reflect.Value, field string) string {
	f := v.FieldByName(field)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Recursively find a field by json:"..." tag.
func findFieldByJSON(t reflect.Type, jsonName string) ([]int, reflect.StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			if idx, fld, ok := findFieldByJSON(f.Type, jsonName); ok {
				return append([]int{i}, idx...), fld, true
			}
		}
		tag := f.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" {
			name = toLowerFirst(f.Name)
		}
		if name == jsonName {
			return []int{i}, f, true
		}
	}
	return nil, reflect.StructField{}, false
}

func findFieldByName(t reflect.Type, name string) ([]int, reflect.StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			if idx, fld, ok := findFieldByName(f.Type, name); ok {
				return append([]int{i}, idx...), fld, true
			}
		}
		if f.Name == name {
			return []int{i}, f, true
		}
	}
	return nil, reflect.StructField{}, false
}

func convertJSONValue(raw any, dest reflect.Type) (reflect.Value, error) {
	// Handle pointers by converting to element type and then taking address.
	if dest.Kind() == reflect.Pointer {
		elem, err := convertJSONValue(raw, dest.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptr := reflect.New(dest.Elem())
		ptr.Elem().Set(elem)
		return ptr, nil
	}

	switch dest.Kind() {
	case reflect.String:
		switch x := raw.(type) {
		case string:
			return reflect.ValueOf(x), nil
		case float64:
			return reflect.ValueOf(strconv.FormatFloat(x, 'f', -1, 64)), nil
		case bool:
			if x {
				return reflect.ValueOf("true"), nil
			}
			return reflect.ValueOf("false"), nil
		default:
			b, _ := json.Marshal(x)
			return reflect.ValueOf(string(b)), nil
		}
	case reflect.Bool:
		if b, ok := raw.(bool); ok {
			return reflect.ValueOf(b), nil
		}
		return reflect.Value{}, errors.New("expected bool")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch x := raw.(type) {
		case float64:
			return reflect.ValueOf(int64(x)).Convert(dest), nil
		case string:
			n, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(dest), nil
		default:
			return reflect.Value{}, errors.New("expected number/string for int")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch x := raw.(type) {
		case float64:
			return reflect.ValueOf(uint64(x)).Convert(dest), nil
		case string:
			n, err := strconv.ParseUint(x, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(dest), nil
		default:
			return reflect.Value{}, errors.New("expected number/string for uint")
		}
	case reflect.Float32, reflect.Float64:
		switch x := raw.(type) {
		case float64:
			return reflect.ValueOf(x).Convert(dest), nil
		case string:
			n, err := strconv.ParseFloat(x, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(dest), nil
		default:
			return reflect.Value{}, errors.New("expected number/string for float")
		}
	case reflect.Slice:
		rawSlice, ok := raw.([]any)
		if !ok {
			return reflect.Value{}, errors.New("expected array")
		}
		slice := reflect.MakeSlice(dest, 0, len(rawSlice))
		for _, elem := range rawSlice {
			cv, err := convertJSONValue(elem, dest.Elem())
			if err != nil {
				return reflect.Value{}, err
			}
			slice = reflect.Append(slice, cv)
		}
		return slice, nil
	case reflect.Struct:
		// Special-case time.Time from RFC3339 strings
		if dest == reflect.TypeOf(time.Time{}) {
			switch x := raw.(type) {
			case string:
				t, err := time.Parse(time.RFC3339, x)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(t), nil
			default:
				return reflect.Value{}, errors.New("expected RFC3339 string for time")
			}
		}
		// Generic: re-marshal & unmarshal
		b, _ := json.Marshal(raw)
		v := reflect.New(dest).Interface()
		if err := json.Unmarshal(b, v); err != nil {
			return reflect.Value{}, err
		}
		return reflect.Indirect(reflect.ValueOf(v)), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported kind %s", dest.Kind())
	}
}

// slug helpers: derive from the first field tagged cms:"slug"
func findSlugSourceValue(v reflect.Value, t reflect.Type) string {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			if s := findSlugSourceValue(v.Field(i), f.Type); s != "" {
				return s
			}
		}
		if f.Tag.Get("cms") == "slug" {
			fv := v.Field(i)
			if fv.Kind() == reflect.String {
				return fv.String()
			}
		}
	}
	return ""
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	out = strings.Join(strings.FieldsFunc(out, func(r rune) bool { return r == '-' }), "-")
	return out
}

//
// ===== 5) Middleware generator (metaprogramming) =====
//

// makeBindAndValidateMiddleware parses JSON into a new instance of reg.Type,
// applies `required:"true"` tag validation, auto-fills slug if empty,
// and puts the instance into context under "cms.obj".
func makeBindAndValidateMiddleware(reg typeReg, store *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		objPtr := reflect.New(reg.Type).Interface()
		if err := c.ShouldBindJSON(objPtr); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := reflect.Indirect(reflect.ValueOf(objPtr))

		// required tag validation
		if err := validateRequired(v, reg.Type); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// timestamps
		now := time.Now().UTC()
		setField(v, "CreatedAt", now)
		setField(v, "UpdatedAt", now)

		// slug
		if getStringField(v, "Slug") == "" {
			if src := findSlugSourceValue(v, reg.Type); src != "" {
				s := slugify(src)
				if s != "" && store.slugExists(reg.Name, s) {
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "slug already exists"})
					return
				}
				setField(v, "Slug", s)
			}
		}

		c.Set("cms.obj", objPtr)
		c.Next()
	}
}

func validateRequired(v reflect.Value, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			if err := validateRequired(fv, f.Type); err != nil {
				return err
			}
		}
		if f.Tag.Get("required") == "true" {
			if isZero(fv) {
				jsonName := f.Tag.Get("json")
				name := strings.Split(jsonName, ",")[0]
				if name == "" {
					name = toLowerFirst(f.Name)
				}
				return fmt.Errorf("field %q is required", name)
			}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// ===== 6) Route generation (metaprogramming) =====
type fieldInfo struct {
	Name     string `json:"name"`
	JSON     string `json:"json"`
	Required bool   `json:"required"`
	IsSlug   bool   `json:"isSlugSource"`
	Kind     string `json:"kind"`
}

func mountRoutes(r *gin.Engine, store *Store) {
	api := r.Group("/api")

	// Introspection: list registered types + their fields
	api.GET("/types", func(c *gin.Context) {
		registryMu.RLock()
		defer registryMu.RUnlock()

		payload := map[string]any{}
		for name, reg := range registry {
			var fields []fieldInfo
			collectFields(reg.Type, &fields)
			payload[name] = fields
		}
		c.JSON(http.StatusOK, payload)
	})

	// For each registered type -> auto CRUD
	registryMu.RLock()
	defer registryMu.RUnlock()

	for name, reg := range registry {

		g := api.Group("/" + name)

		// CREATE
		g.POST("", makeBindAndValidateMiddleware(reg, store), func(c *gin.Context) {
			objPtr, _ := c.Get("cms.obj")
			saved, err := store.Create(name, objPtr)
			if err != nil {
				status := http.StatusInternalServerError
				if strings.Contains(err.Error(), "slug already exists") {
					status = http.StatusConflict
				}
				c.JSON(status, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, saved)
		})

		// LIST
		g.GET("", func(c *gin.Context) {
			list := store.List(name)
			c.JSON(http.StatusOK, list)
		})

		// GET by ID or slug (auto-detect)
		g.GET("/:idOrSlug", func(c *gin.Context) {
			idOrSlug := c.Param("idOrSlug")
			if obj, ok := store.Get(name, idOrSlug); ok {
				c.JSON(http.StatusOK, obj)
				return
			}
			if obj, ok := store.GetBySlug(name, idOrSlug); ok {
				c.JSON(http.StatusOK, obj)
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})

		// PATCH (partial update)
		g.PATCH("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var patch map[string]any
			if err := c.ShouldBindJSON(&patch); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			obj, err := store.Update(name, id, patch, reg)
			if err != nil {
				status := http.StatusInternalServerError
				if strings.Contains(err.Error(), "not found") {
					status = http.StatusNotFound
				}
				if strings.Contains(err.Error(), "slug already exists") {
					status = http.StatusConflict
				}
				c.JSON(status, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, obj)
		})

		// DELETE
		g.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			ok := store.Delete(name, id)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.Status(http.StatusNoContent)
		})
	}
}

func collectFields(t reflect.Type, out *[]fieldInfo) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			collectFields(f.Type, out)
			continue
		}
		jsonName := strings.Split(f.Tag.Get("json"), ",")[0]
		if jsonName == "" {
			jsonName = toLowerFirst(f.Name)
		}
		*out = append(*out, fieldInfo{
			Name:     f.Name,
			JSON:     jsonName,
			Required: f.Tag.Get("required") == "true",
			IsSlug:   f.Tag.Get("cms") == "slug",
			Kind:     f.Type.Kind().String(),
		})
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	store := NewStore()

	// Register any number of content types; routes are generated automatically.
	RegisterContent[Article]("articles")
	RegisterContent[Event]("events")

	mountRoutes(r, store)

	fmt.Println("Dynamic CMS running on http://localhost:8080")
	_ = r.Run(":8080")
}
