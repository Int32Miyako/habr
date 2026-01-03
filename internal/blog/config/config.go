package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		*Database
		*HTTPServer
		*AuthGRPC
	}

	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}

	HTTPServer struct {
		Port        string
		Timeout     time.Duration
		IdleTimeout time.Duration
	}

	AuthGRPC struct {
		Port string
	}
)

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	timeout, err := strconv.Atoi(os.Getenv("BLOG_HTTP_TIMEOUT"))
	if err != nil {
		log.Fatal("Error parsing BLOG_HTTP_TIMEOUT")
	}
	idleTimeout, err := strconv.Atoi(os.Getenv("BLOG_HTTP_IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal("Error parsing BLOG_HTTP_IDLE_TIMEOUT")
	}

	return &Config{
		Database: &Database{
			Host:     os.Getenv("BLOG_DB_HOST"),
			Port:     os.Getenv("BLOG_DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("BLOG_DB_NAME"),
		},

		HTTPServer: &HTTPServer{
			Port:        os.Getenv("BLOG_HTTP_PORT"),
			Timeout:     time.Duration(timeout),
			IdleTimeout: time.Duration(idleTimeout),
		},

		AuthGRPC: &AuthGRPC{
			Port: os.Getenv("AUTH_GRPC_PORT"),
		},
	}
}
