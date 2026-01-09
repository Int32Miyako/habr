package services

import (
	"habr/internal/notification/core/interfaces/repositories"
	"habr/internal/notification/core/interfaces/services"
)

func NewEmailService(repo repositories.EmailRepository) services.EmailService {
	return nil
}

func SendEmail(to string, subject string, body string) error {
	return nil
}
