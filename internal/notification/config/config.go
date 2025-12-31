package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	Database   *Database
	GRPC       *GRPCServer
	AuthClient *AuthClient
}

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

type GRPCServer struct {
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type AuthClient struct {
	Host string
	Port string
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	env := os.Getenv("NOTIFICATION_ENV")
	if env == "" {
		log.Fatal("NOTIFICATION_ENV must be set")
	}

	// Database config
	username := os.Getenv("NOTIFICATION_DB_USER")
	if username == "" {
		log.Fatal("NOTIFICATION_DB_USER must be set")
	}

	password := os.Getenv("NOTIFICATION_DB_PASSWORD")
	if password == "" {
		log.Fatalf("NOTIFICATION_DB_PASSWORD must be set")
	}

	dbHost := os.Getenv("NOTIFICATION_DB_HOST")
	if dbHost == "" {
		log.Fatal("NOTIFICATION_DB_HOST must be set")
	}

	dbPort := os.Getenv("NOTIFICATION_DB_PORT")
	if dbPort == "" {
		log.Fatal("NOTIFICATION_DB_PORT must be set")
	}

	dbName := os.Getenv("NOTIFICATION_DB_NAME")
	if dbName == "" {
		log.Fatalf("NOTIFICATION_DB_NAME must be set")
	}

	// Grpc notification server config
	grpcPort := os.Getenv("NOTIFICATION_GRPC_SERVER_PORT")
	if grpcPort == "" {
		log.Fatal("NOTIFICATION_GRPC_PORT must be set")
	}

	timeoutEnv := os.Getenv("NOTIFICATION_GRPC_TIMEOUT")
	if timeoutEnv == "" {
		log.Fatal("NOTIFICATION_GRPC_TIMEOUT must be set")
	}

	timeout, err := strconv.Atoi(timeoutEnv)
	if err != nil || timeout <= 0 {
		log.Fatal("NOTIFICATION_GRPC_TIMEOUT must be positive integer")
	}

	idleTimeoutEnv := os.Getenv("NOTIFICATION_GRPC_IDLE_TIMEOUT")
	if idleTimeoutEnv == "" {
		log.Fatal("NOTIFICATION_GRPC_IDLE_TIMEOUT must be set")
	}

	idleTimeout, err := strconv.Atoi(idleTimeoutEnv)
	if err != nil || idleTimeout <= 0 {
		log.Fatal("NOTIFICATION_GRPC_IDLE_TIMEOUT must be positive integer")
	}

	// Auth client config
	grpcAuthClientHost := os.Getenv("NOTIFICATION_GRPC_AUTH_CLIENT_HOST")
	if grpcAuthClientHost == "" {
		log.Fatal("NOTIFICATION_GRPC_AUTH_CLIENT_HOST must be set")
	}

	grpcAuthClientPort := os.Getenv("NOTIFICATION_GRPC_AUTH_CLIENT_PORT")
	if grpcAuthClientPort == "" {
		log.Fatal("NOTIFICATION_GRPC_AUTH_CLIENT_PORT must be set")
	}

	return &Config{
		Env: env,
		Database: &Database{
			Username: username,
			Password: password,
			Host:     dbHost,
			Port:     dbPort,
			DBName:   dbName,
		},
		GRPC: &GRPCServer{
			Port:        grpcPort,
			Timeout:     time.Duration(timeout) * time.Second,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
		},
		AuthClient: &AuthClient{
			Host: grpcAuthClientHost,
			Port: grpcAuthClientPort,
		},
	}
}
