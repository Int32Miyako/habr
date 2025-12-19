package http_server

import (
	"habr/internal/blog/core/blog"
	"habr/internal/blog/grpc/client"
	handlers "habr/internal/blog/http-server/handlers/blog"
	"habr/internal/blog/http-server/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(blogService *blog.Service, authClient *client.AuthClient) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/blogs", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authClient))

		r.Get("/", handlers.GetAllBlogs(blogService))

		r.Get("/{id}", handlers.GetBlogByID(blogService))

		r.Post("/", handlers.CreateBlog(blogService))

		r.Put("/{id}", handlers.UpdateBlog(blogService))

		r.Delete("/{id}", handlers.DeleteBlog(blogService))
	})

	return r
}
