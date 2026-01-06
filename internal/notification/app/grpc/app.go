package grpc

import (
	"context"
	"fmt"
	grpcserver "habr/internal/notification/app/grpc/server"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer   *grpc.Server
	emailService services.EmailService
	log          *slog.Logger
	cfg          *config.Config
}

func New(log *slog.Logger, cfg *config.Config, emailService services.EmailService) *App {
	gRPCServer := grpc.NewServer()

	grpcserver.Register(gRPCServer, emailService, log)

	return &App{emailService: emailService, log: log, cfg: cfg, gRPCServer: gRPCServer}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", ":"+app.cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err = app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) {
	const op = "notification grpcapp.Stop"

	app.log.
		Info("stopping notification gRPC server",
			slog.String("op", op),
			slog.String("address", app.cfg.GRPC.Port))

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
