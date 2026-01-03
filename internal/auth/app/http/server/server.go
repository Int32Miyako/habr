package server

import (
	"habr/internal/auth/app/http/router"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"net/http"
	"time"
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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		HTTPServer: srv,
	}
}
