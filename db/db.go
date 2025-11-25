package db

import (
	"context"
	"fmt"
	"habr/internal/blog/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	HOST = "localhost"
	PORT = 5432
)

type Database struct {
	Pool *pgxpool.Pool
}

func Initialize(ctx context.Context, cfg *config.Database) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, cfg.Username, cfg.Password, cfg.DBName)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	// Проверяем подключение
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	_, err = pool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS blogs (
				id SERIAL PRIMARY KEY,
				blog_name TEXT);

             CREATE INDEX IF NOT EXISTS idx_blog_name ON blogs(blog_name);
		`)
	if err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}
