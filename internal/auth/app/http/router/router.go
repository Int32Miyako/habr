package router

import (
	authHandlers "habr/internal/auth/app/http/handlers/auth"
	"habr/internal/auth/app/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(userService *services.UserService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandlers.RegisterUser(userService))

		r.Post("/login", authHandlers.LoginUser(userService))
	})

	return r
}
