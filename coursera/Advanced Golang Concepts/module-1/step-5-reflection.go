package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ValidateStruct(s any) map[string][]string {
	errs := map[string][]string{}
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errs
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		tag := sf.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			switch {
			case rule == "required":
				if isZero(fv) {
					errs[sf.Name] = append(errs[sf.Name], "required")
				}

			case strings.HasPrefix(rule, "min_val="):
				minStr := strings.TrimPrefix(rule, "min_val=")
				min_val, _ := strconv.Atoi(minStr)
				switch fv.Kind() {
				case reflect.String:
					if fv.Len() < min_val {
						errs[sf.Name] = append(errs[sf.Name], fmt.Sprintf("min_val length %d", min_val))
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if fv.Int() < int64(min_val) {
						errs[sf.Name] = append(errs[sf.Name], fmt.Sprintf("min_val value %d", min_val))
					}
				default:
					panic("unhandled default case")
				}
			}
		}
	}
	return errs
}

func isZero(v reflect.Value) bool {
	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	return v.IsZero()
}

func main() {
	type SignUp struct {
		Email string `validate:"required"`
		Name  string `validate:"required,min=3"`
		Age   int    `validate:"min=18"`
	}

	form := SignUp{Email: "x@y.z", Name: "Al", Age: 17}
	validateStruct := ValidateStruct(form)
	fmt.Println(validateStruct)
	// map[Age:[min value 18] Name:[min length 3]]

}
