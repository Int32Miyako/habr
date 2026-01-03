package grpc

import (
	"fmt"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	grpcserver "habr/internal/auth/grpc/server"
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

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", ":"+app.cfg.HTTPServer.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err = app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("address", app.cfg.HTTPServer.Port))

	app.gRPCServer.GracefulStop()
}
