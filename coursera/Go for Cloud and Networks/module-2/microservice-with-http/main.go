package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	mu   sync.RWMutex
	data map[int64]User
}

func NewUserStore() *UserStore {
	return &UserStore{data: make(map[int64]User)}
}

func (s *UserStore) Create(user User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = time.Now().UnixNano() + int64(rand.Intn(1000))
	s.data[user.ID] = user
	return user
}

func (s *UserStore) Get(id int64) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.data[id]
	return user, ok
}

func (s *UserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}

func (s *UserStore) Update(id int64, user User) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return User{}, false
	}
	user.ID = id
	s.data[id] = user
	return user, true
}

func (s *UserStore) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return false
	}
	delete(s.data, id)
	return true
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

type writeWrap struct {
	http.ResponseWriter
	code int
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &writeWrap{ResponseWriter: w, code: http.StatusOK}
			next.ServeHTTP(ww, r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.code,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote", r.RemoteAddr)
		})
	}
}

func (w *writeWrap) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func recoverMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic", "err", rec)
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func authMiddleware(expectedBearer string) func(http.Handler) http.Handler {
	// very simple example: require "Authorization: Bearer <token>"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expectedBearer == "" {
				next.ServeHTTP(w, r) // auth disabled
				return
			}
			got := r.Header.Get("Authorization")
			if !strings.HasPrefix(got, "Bearer ") || strings.TrimPrefix(got, "Bearer ") != expectedBearer {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type Server struct {
	log   *slog.Logger
	store *UserStore
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.store.List())
	case http.MethodPost:
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := readJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name and email are required"})
			return
		}
		u := s.store.Create(User{Name: req.Name, Email: req.Email})
		writeJSON(w, http.StatusCreated, u)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) userByID(w http.ResponseWriter, r *http.Request) {
	// Expect /users/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		if u, ok := s.store.Get(id); ok {
			writeJSON(w, http.StatusOK, u)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
	case http.MethodPut:
		var req User
		if err := readJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name and email are required"})
			return
		}
		if u, ok := s.store.Update(id, req); ok {
			writeJSON(w, http.StatusOK, u)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
	case http.MethodDelete:
		if s.store.Delete(id) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func newMux(s *Server, logger *slog.Logger, bearer string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /users", s.users)
	mux.HandleFunc("POST /users", s.users)
	mux.HandleFunc("/users/", s.userByID) // GET|PUT|DELETE /users/{id}

	// chain: recover -> logging -> auth -> mux
	return recoverMiddleware(logger)(
		loggingMiddleware(logger)(
			authMiddleware(bearer)(mux),
		),
	)
}

func main() {
	// Config via env (minimal)
	port := getenvDefault("PORT", "8080")
	bearer := os.Getenv("AUTH_BEARER") // if empty, auth disabled

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	srv := &Server{
		log:   logger,
		store: NewUserStore(),
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           newMux(srv, logger, bearer),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start
	go func() {
		logger.Info("server starting", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("server shutting down...")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}
	logger.Info("bye")
}

func getenvDefault(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
