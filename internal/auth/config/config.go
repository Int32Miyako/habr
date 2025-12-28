package config

import (
	"habr/internal/auth/core/constants"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		*Database
		*HTTPServer
		*JWT
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

type JWT struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	timeout, err := strconv.Atoi(os.Getenv("AUTH_GRPC_TIMEOUT"))
	if err != nil {
		timeout = constants.DefaultGRPCTimeoutSeconds
	}

	accessTokenDuration, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_DURATION_MINUTES"))
	if err != nil {
		accessTokenDuration = constants.DefaultAccessTokenDurationMinutes
	}

	refreshTokenDuration, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_DURATION_DAYS"))
	if err != nil {
		refreshTokenDuration = constants.DefaultRefreshTokenDurationDays
	}

	return &Config{
		Database: &Database{
			Host:     os.Getenv("AUTH_DB_HOST"),
			Port:     os.Getenv("AUTH_DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("AUTH_DB_NAME"),
		},

		HTTPServer: &HTTPServer{
			Address:     os.Getenv("AUTH_GRPC_ADDRESS"),
			Timeout:     time.Duration(timeout),
			IdleTimeout: time.Duration(timeout),
		},

		JWT: &JWT{
			SecretKey:            os.Getenv("JWT_SECRET_KEY"),
			AccessTokenDuration:  time.Duration(accessTokenDuration) * time.Minute,
			RefreshTokenDuration: time.Duration(refreshTokenDuration) * 24 * time.Hour,
		},
	}

}
