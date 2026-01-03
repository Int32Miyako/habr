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

// Start запускает Kafka consumer и gRPC сервер.
func (app *App) Start(ctx context.Context) error {
	// Запускаем Kafka consumer
	err := app.KafkaApp.RegistrationConsumer.Start(ctx)
	if err != nil {
		return fmt.Errorf("kafka consumer start: %w", err)
	}

	// Запускаем gRPC сервер
	err = app.GRPCApp.Run()
	if err != nil {
		return fmt.Errorf("kafka consumer start: %w", err)
	}

	return nil
}

// Stop останавливает Kafka consumer и gRPC сервер.
func (app *App) Stop() error {
	if app.KafkaApp.RegistrationConsumer != nil {
		if err := app.KafkaApp.Stop(); err != nil {
			return fmt.Errorf("app stop: %w", err)
		}
	}

	if app.GRPCApp != nil {
		app.GRPCApp.Stop()
	}

	return nil
}
