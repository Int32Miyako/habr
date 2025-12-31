package repositories

import "habr/internal/notification/core/models"

type EmailRepository interface {
	Create(email *models.Email) (int64, error)
	GetAll() ([]*models.Email, error)
}
