package models

import "time"

type User struct {
	ID           int64     `dbcodes:"id"`
	Username     string    `dbcodes:"username"`
	PasswordHash string    `dbcodes:"password_hash"`
	Email        string    `dbcodes:"email"`
	CreatedAt    time.Time `dbcodes:"created_at"`
	UpdatedAt    time.Time `dbcodes:"updated_at"`
}

type RefreshToken struct {
	ID        int64     `dbcodes:"id"`
	UserID    int64     `dbcodes:"user_id"`
	TokenHash string    `dbcodes:"token_hash"`
	ExpiresAt time.Time `dbcodes:"expires_at"`
	CreatedAt time.Time `dbcodes:"created_at"`
}
