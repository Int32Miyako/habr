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
	}

	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}

	HTTPServer struct {
		Address     string
		Timeout     time.Duration
		IdleTimeout time.Duration
	}
)

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	timeout, err := strconv.Atoi(os.Getenv("AUTH_GRPC_TIMEOUT"))
	if err != nil {
		log.Fatal("Error parsing HTTP_TIMEOUT")
	}

	return &Config{
		Database: &Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("AUTH_DB_NAME"),
		},

		HTTPServer: &HTTPServer{
			Address:     os.Getenv("AUTH_GRPC_ADDRESS"),
			Timeout:     time.Duration(timeout),
			IdleTimeout: time.Duration(timeout),
		},
	}

}
