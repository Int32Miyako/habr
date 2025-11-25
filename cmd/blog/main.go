package main

import (
	"context"
	"habr/db"
	"habr/internal/blog/config"
	"habr/internal/blog/core/blog"
	"habr/internal/blog/http-server"
	"log"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	database, err := db.Initialize(ctx, cfg.Database)
	if err != nil {
		panic(err)
	}

	blogRepository := blog.NewRepository(database.Pool)
	blogService := blog.NewService(blogRepository)

	router := http_server.NewRouter(blogService)

	log.Println("listening on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
