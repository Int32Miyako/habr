package app

import (
	"habr/internal/auth/config"
	"log/slog"
	"net"
)

type App struct {
	// GRPCServer *grpcapp.App
}

func New() *App {
	return &App{}
}

func (app *App) Start(cfg *config.Config, log *slog.Logger) {
	l, err := net.Listen("tcp", cfg.Port)
	defer l.Close()
	if err != nil {
		log.Error("Error starting http server", "error", err)
	}

	log.Info("Starting http server", "port", cfg.Port)
}
