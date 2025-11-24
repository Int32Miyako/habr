package main

import (
	"context"
	"fmt"
	"habr/db"
	"habr/internal/config"
	"habr/internal/core/blog"
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

	var id int64
	id, err = blogService.CreateBlog(ctx, "yes")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(id)

	blogs, err := blogService.GetBlogs(ctx)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(blogs)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
