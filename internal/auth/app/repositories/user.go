package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, username, passwordHash, tokenHash string, expiresAt time.Time) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int64
	err = tx.QueryRow(ctx, query, username, passwordHash, email).Scan(&userId)
	if err != nil {
		return -1, err
	}

	query = `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, query, userId, tokenHash, expiresAt)
	if err != nil {
		return -1, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (int64, string, error) {
	query := `SELECT id, password_hash FROM users WHERE email = $1`

	var userID int64
	var passwordHash string
	err := r.pool.QueryRow(ctx, query, email).Scan(&userID, &passwordHash)
	if err != nil {
		return -1, "", err
	}

	return userID, passwordHash, nil
}
