package http

import (
	"habr/internal/blog/core/blog"
	"habr/internal/blog/grpc/client"
	blogHandlers "habr/internal/blog/http/handlers/blog"
	"habr/internal/blog/http/middlewares"
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
	r.Use(middlewares.AuthMiddleware(authClient))

	r.Route("/blogs", func(r chi.Router) {

		r.Get("/", blogHandlers.GetAllBlogs(blogService))

		r.Get("/{id}", blogHandlers.GetBlogByID(blogService))

		r.Post("/", blogHandlers.CreateBlog(blogService))

		r.Put("/{id}", blogHandlers.UpdateBlog(blogService))

		r.Delete("/{id}", blogHandlers.DeleteBlog(blogService))

	})

	return r
}
