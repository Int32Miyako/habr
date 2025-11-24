package main

import (
	"context"
	"habr/db"
	"habr/internal/config"
	"habr/internal/core/blog"
	"log"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	database, err := db.Initialize(ctx, cfg.Database)
	if err != nil {
		panic(err)
	}

	blogRepository := blog.NewRepository(database.Pool)
	err = blogRepository.CreateBlog(ctx, "ok")
	if err != nil {
		panic(err)
	}

	blogs, err := blogRepository.GetBlogs(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(blogs)
}
