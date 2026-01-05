package grpc

import (
	"context"
	"fmt"
	grpcserver "habr/internal/auth/app/grpc/server"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer  *grpc.Server
	userService *services.UserService
	log         *slog.Logger
	cfg         *config.Config
}

func New(log *slog.Logger, cfg *config.Config, userService *services.UserService) *App {
	gRPCServer := grpc.NewServer()

	grpcserver.Register(gRPCServer, userService, log)

	return &App{userService: userService, log: log, cfg: cfg, gRPCServer: gRPCServer}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", ":"+app.cfg.GRPCServer.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.log.Info("gRPC auth server started", slog.String("addr", l.Addr().String()))

	if err = app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("address", app.cfg.GRPCServer.Port))

	done := make(chan struct{})

	go func() {
		app.gRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		app.log.Info("gRPC server stopped gracefully")

	case <-ctx.Done():
		app.log.Warn("gRPC graceful shutdown timeout exceeded, forcing stop")
		app.gRPCServer.Stop()
	}
}
