package app

import (
	"context"
	grpcapp "habr/internal/auth/app/grpc"
	httpapp "habr/internal/auth/app/http"
	"habr/internal/auth/app/kafka"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/app/kafka/producer/client"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
)

type App struct {
	GRPC     *grpcapp.App
	HTTP     *httpapp.App
	KafkaApp *kafka.App
	log      *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger, userService *services.UserService) *App {
	grpcApp := grpcapp.New(log, cfg, userService)
	httpApp := httpapp.New(log, cfg, userService)

	producerClient, err := client.NewProducerKafkaClient(cfg.Kafka.Brokers, cfg.Kafka.Topic, log)
	if err != nil {
		log.Error("Failed to create kafka registration notifier", "error", err)
	}
	prod := producer.NewRegistrationNotifier(producerClient, log)

	kafkaApp := kafka.New(prod, log)

	return &App{
		GRPC:     grpcApp,
		HTTP:     httpApp,
		KafkaApp: kafkaApp,
		log:      log,
	}
}

func (app *App) Start(serverErrors chan<- error) {
	app.log.Info("Starting HTTP and gRPC servers...")

	go func() {
		if err := app.HTTP.Run(); err != nil {
			serverErrors <- err
		}
	}()

	go func() {
		if err := app.GRPC.Run(); err != nil {
			serverErrors <- err
		}
	}()

	app.log.Info("Servers started in background goroutines")
}

func (app *App) Stop(ctx context.Context) {
	if app.HTTP != nil {
		app.HTTP.Stop(ctx)
	}

	if app.GRPC != nil {
		app.GRPC.Stop(ctx)
	}

	if app.KafkaApp != nil {
		err := app.KafkaApp.Close()
		if err != nil {
			app.log.Error("failed to close kafka prod", slog.String("error", err.Error()))
		}
	}
}
