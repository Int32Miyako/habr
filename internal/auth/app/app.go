package app

import (
	"context"
	grpcapp "habr/internal/auth/app/grpc"
	httpapp "habr/internal/auth/app/http"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
)

type App struct {
	GRPC *grpcapp.App
	HTTP *httpapp.App
}

func New(cfg *config.Config, log *slog.Logger, userService *services.UserService) *App {
	grpcApp := grpcapp.New(log, cfg, userService)
	httpApp := httpapp.New(log, cfg, userService)

	return &App{
		GRPC: grpcApp,
		HTTP: httpApp,
	}
}

func (app *App) Start() {
	go app.HTTP.MustRun()

	app.GRPC.MustRun()
}

func (app *App) Stop(ctx context.Context) {
	if app.HTTP != nil {
		app.HTTP.Stop(ctx)
	}

	if app.GRPC != nil {
		app.GRPC.Stop()
	}
}
