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

func (r *UserRepository) CreateUser(ctx context.Context, email, username, passwordHash, tokenHash string) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING userId`
	var userId int64
	err = tx.QueryRow(ctx, query, username, passwordHash, email).Scan(&userId)
	if err != nil {
		return -1, err
	}

	query = `INSERT INTO refresh_tokens (userId, token_hash) VALUES ($1, $2)`
	tx.QueryRow(ctx, query, userId, tokenHash)

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (r *UserRepository) LoginUser(ctx context.Context, email string, hash string) (int64, error) {
	query := `
				SELECT token_hash FROM refresh_tokens rt
				LEFT JOIN users u ON u.id = rt.user_id
				WHERE u.email = $1
          			AND rt.token_hash = $2
			`

	var userID int64
	err := r.pool.QueryRow(ctx, query, email, hash).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}
