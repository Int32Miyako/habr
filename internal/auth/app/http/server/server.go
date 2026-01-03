package server

import (
	"habr/internal/auth/app/http/router"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"net/http"
)

type Server struct {
	HTTPServer *http.Server
}

func New(cfg *config.Config, userService *services.UserService) *Server {
	r := router.New(userService)
	addr := ":" + cfg.HTTPServer.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout * 2,
		IdleTimeout:  cfg.HTTPServer.Timeout * 15,
	}

	return &Server{
		HTTPServer: srv,
	}
}
