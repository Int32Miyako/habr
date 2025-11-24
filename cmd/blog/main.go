package main

import (
	"context"
	"habr/db"
	"habr/internal/config"
	"habr/internal/core/blog"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	database, err := db.Initialize(ctx, cfg.Database)
	if err != nil {
		panic(err)
	}

	blogRepository := blog.NewRepository(database.Pool)
	err = blogRepository.CreateBlog("ok", ctx)
	if err != nil {
		panic(err)
	}
}
