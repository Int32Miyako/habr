package models

import "time"

type Email struct {
	ID        int64     `db:"id"`
	Address   string    `db:"address"`
	Status    string    `db:"status"` // pending, confirmed, unsubscribed
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
