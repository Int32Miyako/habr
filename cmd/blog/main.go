package main

import (
	"context"
	db "habr/db/blog"
	"habr/internal/blog/config"
	"habr/internal/blog/core/blog"
	"habr/internal/blog/grpc/client"
	"habr/internal/blog/http-server"
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
	authClient, err := client.NewAuthClient(cfg.AuthGRPC.Address)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}
	defer authClient.Close()

	blogRepository := blog.NewRepository(database.Pool)
	blogService := blog.NewService(blogRepository)

	router := http_server.NewRouter(blogService, authClient)

	log.Println("listening on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
