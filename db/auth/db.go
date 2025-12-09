package auth

import (
	"context"
	"fmt"
	"habr/internal/auth/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func Initialize(ctx context.Context, cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	// Проверяем подключение
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Database{Pool: pool}, nil
}
