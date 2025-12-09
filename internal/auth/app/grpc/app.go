package grpc

import (
	"fmt"
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
	return &App{userService: userService, log: log, cfg: cfg, gRPCServer: gRPCServer}
}

func (app *App) Run() {
	const op = "grpcapp.Run"
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.cfg.Port))
	if err != nil {
		panic(err.Error() + op)
	}

	app.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

}
