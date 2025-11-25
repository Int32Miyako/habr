package http_server

import (
	"context"
	"fmt"
	"habr/internal/blog/core/blog"
	"habr/internal/lib/api/http-server"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(ctx context.Context, blogService *blog.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		blogs, err := blogService.GetBlogs(ctx)
		if err != nil {
			log.Println(err)
		}
		http_server.RespJSON(w, blogs)
		fmt.Println(blogs)
	})

	return r
}
