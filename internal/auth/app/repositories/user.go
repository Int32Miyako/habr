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

func (r *UserRepository) CreateUser(ctx context.Context, email, username, passwordHash string) (int64, error) {
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

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (r *UserRepository) CreateRefreshToken(ctx context.Context, userId int64, tokenHash string, expiresAt time.Time) (int64, error) {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) 
				VALUES ($1, $2, $3) RETURNING id`
	var tokenId int64
	err := r.pool.QueryRow(ctx, query, userId, tokenHash, expiresAt).Scan(&tokenId)
	if err != nil {
		return -1, err
	}

	return tokenId, nil
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

func (r *UserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (int64, time.Time, error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`

	var userID int64
	var expiresAt time.Time
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return -1, time.Time{}, err
	}

	return userID, expiresAt, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	return err
}

func (r *UserRepository) GetUserEmailByID(ctx context.Context, userID int64) (string, error) {
	query := `SELECT email FROM users WHERE id = $1`

	var email string
	err := r.pool.QueryRow(ctx, query, userID).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}
