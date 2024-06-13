package main

import (
	"context"
	"fmt"
	"github.com/stashchenko/microservice-example/pkg/postgresutil"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/stashchenko/microservice-example/internal/app/user"
	"github.com/stashchenko/microservice-example/internal/grpc"
	"github.com/stashchenko/microservice-example/internal/grpc/handler"
	"github.com/stashchenko/microservice-example/internal/grpc/health"
	userhandler "github.com/stashchenko/microservice-example/internal/grpc/user"
	"github.com/stashchenko/microservice-example/pkg/config"
	"github.com/stashchenko/microservice-example/pkg/proto"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Updates slogs default instance of slog with our own handler.
	conf := config.NewConfig()

	if err := conf.Load(".env"); err != nil {
		logger.Error("failed to load service config", "error", err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConf := conf.Database
	pgConn, err := postgresutil.Connect(ctx, dbConf.Host, dbConf.Username, dbConf.Password, dbConf.DatabaseName, dbConf.Port)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err.Error())
		os.Exit(1)
	}

	userSvc := user.NewService(user.NewUserRepository(pgConn))
	baseHandler := handler.NewHandler(&handler.Service{User: userSvc})

	grpcServer, err := grpc.NewServer(baseHandler, grpc.WithPort(conf.Server.Port))
	if err != nil {
		logger.Error("failed to init new grpc server", "error", err)

		os.Exit(1)
	}

	proto.RegisterHealthServer(grpcServer.Server(), health.NewHandler())
	proto.RegisterUserServer(grpcServer.Server(), userhandler.NewHandler(baseHandler))

	reflection.Register(grpcServer.Server())

	go func() {
		err := grpcServer.Serve()
		if err != nil {
			logger.Error("error serving grpc server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		logger.Error(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.Error(fmt.Sprintf("ctx.Done: %v", done))
	}
}
