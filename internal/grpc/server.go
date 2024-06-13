package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/stashchenko/microservice-example/internal/grpc/handler"
	"google.golang.org/grpc"
)

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

type Server struct {
	grpc        *grpc.Server
	port        int
	baseHandler *handler.Handler
}

func (s *Server) Server() *grpc.Server {
	return s.grpc
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf(":%d", s.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}
	slog.Info("grpc server listening", "addr", addr)
	return s.grpc.Serve(l)
}

func NewServer(h *handler.Handler, opts ...ServerOption) (*Server, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
		grpc.StreamInterceptor(LoggingStreamInterceptor),
	)

	server := &Server{
		port:        8000, // Default port
		grpc:        grpcServer,
		baseHandler: h,
	}

	// Apply server options
	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Calls the handler
	h, err := handler(ctx, req)

	// Logs the processing time
	slog.Info("Unary Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)

	return h, err
}

func LoggingStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	// Calls the handler
	err := handler(srv, ss)

	// Logs the processing time
	slog.Info("Stream Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)

	return err
}
