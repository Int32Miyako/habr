package config

import (
	"habr/internal/auth/core/constants"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Database   *Database
		HTTPServer *HTTPServer
		JWT        *JWT
		Kafka      *Kafka
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

	Kafka struct {
		Brokers []string
		Topic   string
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

	accessTokenDuration, err := strconv.Atoi(os.Getenv("AUTH_JWT_ACCESS_TOKEN_DURATION_MINUTES"))
	if err != nil {
		accessTokenDuration = constants.DefaultAccessTokenDurationMinutes
	}

	refreshTokenDuration, err := strconv.Atoi(os.Getenv("AUTH_JWT_REFRESH_TOKEN_DURATION_DAYS"))
	if err != nil {
		refreshTokenDuration = constants.DefaultRefreshTokenDurationDays
	}

	// Kafka config
	kafkaBrokersStr := os.Getenv("AUTH_KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		log.Fatal("AUTH_KAFKA_BROKERS must be set")
	}
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	kafkaTopic := os.Getenv("AUTH_KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatal("AUTH_KAFKA_TOPIC must be set")
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
			SecretKey:            os.Getenv("AUTH_JWT_SECRET_KEY"),
			AccessTokenDuration:  time.Duration(accessTokenDuration) * time.Minute,
			RefreshTokenDuration: time.Duration(refreshTokenDuration) * 24 * time.Hour,
		},
		Kafka: &Kafka{
			Brokers: kafkaBrokers,
			Topic:   kafkaTopic,
		},
	}

}
