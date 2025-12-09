package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, username, passwordHash string) (int64, error) {
	query := "INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := r.pool.QueryRow(ctx, query, username, passwordHash, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
