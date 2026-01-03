package http

import (
	"habr/internal/auth/app/http/server"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
	"net/http"
)

type App struct {
	HTTPServer *http.Server
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config, userService *services.UserService) *App {
	srv := server.New(cfg, userService)

	log.Info("HTTP server created", slog.String("addr", srv.HTTPServer.Addr))

	return &App{
		HTTPServer: srv.HTTPServer,
		log:        log,
	}
}
