package app

import (
	"context"
	"fmt"
	"habr/internal/notification/app/grpc"
	"habr/internal/notification/app/kafka"

	"golang.org/x/sync/errgroup"
)

// App представляет собой основное приложение, содержащее gRPC сервер и Kafka consumer.
type App struct {
	GRPCApp  *grpc.App
	KafkaApp *kafka.App
}

// New создает новый экземпляр App с предоставленными gRPC сервером и Kafka consumer.
func New(gRPCApp *grpc.App, kafkaApp *kafka.App) *App {
	return &App{
		GRPCApp:  gRPCApp,
		KafkaApp: kafkaApp,
	}
}

// Start создает контекст с отменой для управления внутренними горутинами
func (app *App) Start(ctx context.Context) error {

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.KafkaApp.Run(gCtx)
	})

	g.Go(func() error {
		return app.GRPCApp.Run()
	})

	return g.Wait()
}

// Stop останавливает Kafka consumer и gRPC сервер.
func (app *App) Stop(shutdownCtx context.Context) error {
	if app.KafkaApp.RegistrationConsumer != nil {
		err := app.KafkaApp.Stop(shutdownCtx)
		if err != nil {
			return fmt.Errorf("app stop: %w", err)
		}
	}

	if app.GRPCApp != nil {
		app.GRPCApp.Stop(shutdownCtx)
	}

	return nil
}
