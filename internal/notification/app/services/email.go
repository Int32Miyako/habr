package services

import (
	"habr/internal/notification/core/interfaces/repositories"
	"habr/internal/notification/core/interfaces/services"
)

func NewEmailService(repo repositories.EmailRepository) services.EmailService {
	return nil
}
