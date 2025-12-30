package config

import "github.com/joho/godotenv"

type Config struct {
	DB   *Database
	GRPC *GRPCServer
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type GRPCServer struct {
	Address     string
	Timeout     int
	IdleTimeout int
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		DB:   &Database{},
		GRPC: &GRPCServer{},
	}
}
