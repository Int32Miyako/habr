package main

import (
	"context"
	db "habr/db/blog"
	"habr/internal/blog/config"
	"habr/internal/blog/core/blog"
	"habr/internal/blog/grpc/client"
	httpserver "habr/internal/blog/http"
	"log"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	database, err := db.Initialize(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// Инициализация gRPC клиента для сервиса аутентификации
	authClient, err := client.NewAuthClient(cfg.AuthGRPC.Port)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	defer func() {
		if err = authClient.Close(); err != nil {
			log.Printf("failed to close auth client: %v", err)
		}
	}()

	blogRepository := blog.NewRepository(database.Pool)
	blogService := blog.NewService(blogRepository)

	router := httpserver.NewRouter(blogService, authClient)

	log.Printf("listening on :%s", cfg.HTTPServer.Port)

	err = http.ListenAndServe(":"+cfg.HTTPServer.Port, router)
	if err != nil {
		panic(err)
	}
}
