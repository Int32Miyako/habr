package app

import (
	"habr/internal/notification/app/grpc"

	"github.com/IBM/sarama"
)

type App struct {
	GRPCServer *grpc.App
	Consumer   sarama.Consumer
}

func New(gRPCServer *grpc.App, consumer sarama.Consumer) *App {
	return &App{
		GRPCServer: gRPCServer,
		Consumer:   consumer,
	}
}

func (app *App) Start() error {
	err := app.GRPCServer.Run()
	return err
}

func (app *App) Stop() {
	if app.GRPCServer != nil {
		app.GRPCServer.Stop()
	}
}
