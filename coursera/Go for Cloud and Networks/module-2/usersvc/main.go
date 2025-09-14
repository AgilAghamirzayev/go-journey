package main

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"log"
//	"log/slog"
//	"math/rand/v2"
//	"net"
//	"os"
//	"os/signal"
//	"strings"
//	"sync"
//	"syscall"
//	"time"
//
//	userv1 "usersvc/gen/user/v1"
//
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/credentials"
//	"google.golang.org/grpc/credentials/insecure"
//	grpc_health "google.golang.org/grpc/health"
//	"google.golang.org/grpc/health/grpc_health_v1"
//	"google.golang.org/grpc/metadata"
//	"google.golang.org/grpc/reflection"
//	"google.golang.org/grpc/status"
//)
//
//// ---------- In-memory store ----------
//type store struct {
//	mu   sync.RWMutex
//	data map[int64]*userv1.User
//}
//
//func newStore() *store { return &store{data: make(map[int64]*userv1.User)} }
//func (s *store) create(u *userv1.User) *userv1.User {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	u.Id = time.Now().UnixNano() + int64(rand.IntN(1000))
//	cp := *u
//	s.data[cp.Id] = &cp
//	return &cp
//}
//func (s *store) get(id int64) (*userv1.User, bool) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//	u, ok := s.data[id]
//	if !ok {
//		return nil, false
//	}
//	cp := *u
//	return &cp, true
//}
//func (s *store) list() []*userv1.User {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//	out := make([]*userv1.User, 0, len(s.data))
//	for _, u := range s.data {
//		cp := *u
//		out = append(out, &cp)
//	}
//	return out
//}
//func (s *store) update(id int64, in *userv1.User) (*userv1.User, bool) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	if _, ok := s.data[id]; !ok {
//		return nil, false
//	}
//	cp := *in
//	cp.Id = id
//	s.data[id] = &cp
//	return &cp, true
//}
//func (s *store) delete(id int64) bool {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	if _, ok := s.data[id]; !ok {
//		return false
//	}
//	delete(s.data, id)
//	return true
//}
//
//// ---------- Service impl ----------
//type userService struct {
//	userv1.UnimplementedUserServiceServer
//	log   *slog.Logger
//	store *store
//}
//
//func (s *userService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error) {
//	if req.GetId() <= 0 {
//		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
//	}
//	u, ok := s.store.get(req.GetId())
//	if !ok {
//		return nil, status.Error(codes.NotFound, "user not found")
//	}
//	return u, nil
//}
//
//func (s *userService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
//	name := strings.TrimSpace(req.GetName())
//	email := strings.TrimSpace(req.GetEmail())
//	if name == "" || email == "" {
//		return nil, status.Error(codes.InvalidArgument, "name and email are required")
//	}
//	u := s.store.create(&userv1.User{Name: name, Email: email})
//	return u, nil
//}
//
//func (s *userService) ListUsers(ctx context.Context, _ *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
//	return &userv1.ListUsersResponse{Users: s.store.list()}, nil
//}
//
//func (s *userService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
//	if req.GetId() <= 0 {
//		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
//	}
//	name := strings.TrimSpace(req.GetName())
//	email := strings.TrimSpace(req.GetEmail())
//	if name == "" || email == "" {
//		return nil, status.Error(codes.InvalidArgument, "name and email are required")
//	}
//	u, ok := s.store.update(req.GetId(), &userv1.User{Name: name, Email: email})
//	if !ok {
//		return nil, status.Error(codes.NotFound, "user not found")
//	}
//	return u, nil
//}
//
//func (s *userService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
//	if req.GetId() <= 0 {
//		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
//	}
//	ok := s.store.delete(req.GetId())
//	return &userv1.DeleteUserResponse{Ok: ok}, nil
//}
//
//// ---------- Interceptors (logging, panic-recover, bearer auth) ----------
//func loggingUnary(logger *slog.Logger) grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
//		start := time.Now()
//		resp, err := handler(ctx, req)
//		code := status.Code(err)
//		logger.Info("grpc",
//			"method", info.FullMethod,
//			"code", code.String(),
//			"ms", time.Since(start).Milliseconds(),
//		)
//		return resp, err
//	}
//}
//func recoverUnary(logger *slog.Logger) grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
//		defer func() {
//			if r := recover(); r != nil {
//				logger.Error("panic", "err", r)
//				err = status.Error(codes.Internal, "internal server error")
//			}
//		}()
//		return handler(ctx, req)
//	}
//}
//func authUnary(expectedBearer string) grpc.UnaryServerInterceptor {
//	// very simple: require metadata "authorization: Bearer <token>" if configured
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
//		if expectedBearer == "" {
//			return handler(ctx, req)
//		}
//		md, ok := metadata.FromIncomingContext(ctx)
//		if !ok {
//			return nil, status.Error(codes.Unauthenticated, "missing metadata")
//		}
//		values := md.Get("authorization")
//		if len(values) == 0 || !strings.HasPrefix(strings.ToLower(values[0]), "bearer ") {
//			return nil, status.Error(codes.Unauthenticated, "missing bearer token")
//		}
//		got := strings.TrimSpace(values[0][len("bearer "):])
//		if got != expectedBearer {
//			return nil, status.Error(codes.Unauthenticated, "invalid token")
//		}
//		return handler(ctx, req)
//	}
//}
//
//// ---------- Main ----------
//func main() {
//	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
//
//	port := getenv("PORT", "50051")
//	addr := ":" + port
//	authToken := os.Getenv("AUTH_BEARER")
//	enableReflection := getenv("REFLECTION", "1") == "1"
//
//	var creds grpc.ServerOption
//	cert := os.Getenv("TLS_CERT_FILE")
//	key := os.Getenv("TLS_KEY_FILE")
//	if cert != "" && key != "" {
//		tlsCreds, err := credentials.NewServerTLSFromFile(cert, key)
//		if err != nil {
//			log.Fatalf("load tls creds: %v", err)
//		}
//		creds = grpc.Creds(tlsCreds)
//	} else {
//		creds = grpc.Creds(insecure.NewCredentials()) // for local/dev; remove in prod
//	}
//
//	s := grpc.NewServer(
//		creds,
//		grpc.ChainUnaryInterceptor(
//			recoverUnary(logger),
//			loggingUnary(logger),
//			authUnary(authToken),
//		),
//	)
//
//	// Services
//	us := &userService{log: logger, store: newStore()}
//	userv1.RegisterUserServiceServer(s, us)
//
//	// Health + (opt) reflection
//	hs := grpc_health.NewServer()
//	grpc_health_v1.RegisterHealthServer(s, hs)
//	if enableReflection {
//		reflection.Register(s)
//	}
//
//	lis, err := net.Listen("tcp", addr)
//	if err != nil {
//		log.Fatalf("listen: %v", err)
//	}
//	logger.Info("gRPC server starting", "addr", addr)
//
//	// serve in background
//	errCh := make(chan error, 1)
//	go func() {
//		if e := s.Serve(lis); !errors.Is(e, grpc.ErrServerStopped) {
//			errCh <- e
//		}
//	}()
//
//	// graceful shutdown
//	sigCh := make(chan os.Signal, 1)
//	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
//	select {
//	case sig := <-sigCh:
//		logger.Info("signal received; stopping", "sig", sig.String())
//		s.GracefulStop()
//	case e := <-errCh:
//		logger.Error("server error", "err", e)
//		os.Exit(1)
//	}
//	logger.Info("bye")
//}
//
//func getenv(k, d string) string {
//	if v := os.Getenv(k); v != "" {
//		return v
//	}
//	return d
//}
