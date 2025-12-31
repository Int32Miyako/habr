package common

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database хранит пул соединений
type Database struct {
	Pool *pgxpool.Pool
}

// DBConfig минимальная структура конфигурации БД
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

// Initialize создаёт пул соединений и проверяет подключение
func Initialize(ctx context.Context, cfg *DBConfig) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName,
	)

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

// Close закрывает пул соединений
func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
