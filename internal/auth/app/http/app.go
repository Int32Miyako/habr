package http

import (
	"context"
	"errors"
	"fmt"
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
	srv := server.New(cfg, userService).HTTPServer

	return &App{
		HTTPServer: srv,
		log:        log,
	}
}

func (app *App) Run() error {
	app.log.Info("HTTP auth server started", slog.String("addr", app.HTTPServer.Addr))

	if err := app.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server start: %w", err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) {
	const op = "httpapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping HTTP server", slog.String("addr", app.HTTPServer.Addr))

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		app.log.Error("failed to stop HTTP server", slog.String("error", err.Error()))
	}

	app.log.Info("HTTP server stopped")
}
