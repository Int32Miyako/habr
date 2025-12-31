package repositories

import (
	"habr/internal/notification/core/interfaces/repositories"
	"habr/internal/notification/core/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailRepository struct {
	pool *pgxpool.Pool
}

func (e EmailRepository) Create(email *models.Email) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmailRepository) GetAll() ([]*models.Email, error) {
	//TODO implement me
	panic("implement me")
}

func NewEmailRepository(pool *pgxpool.Pool) repositories.EmailRepository {
	return EmailRepository{pool: pool}
}
