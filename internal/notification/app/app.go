package app

import (
	"context"
	"fmt"
	"habr/internal/notification/app/grpc"
	"habr/internal/notification/app/kafka"
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

// Start запускает Kafka consumer и gRPC сервер в отдельных горутинах.
func (app *App) Start(ctx context.Context) error {
	errChan := make(chan error, 2)

	go func() {
		errChan <- app.KafkaApp.Run(ctx)
	}()

	go func() {
		errChan <- app.GRPCApp.Run()
	}()

	return <-errChan
}

// Stop останавливает Kafka consumer и gRPC сервер.
func (app *App) Stop(ctx context.Context) error {
	if app.KafkaApp.RegistrationConsumer != nil {
		err := app.KafkaApp.Stop(ctx)
		if err != nil {
			return fmt.Errorf("app stop: %w", err)
		}
	}

	if app.GRPCApp != nil {
		app.GRPCApp.Stop(ctx)
	}

	return nil
}
