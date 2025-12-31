package services

type EmailService interface {
	Subscribe(email string) (int64, error)
	ListSubscribers() ([]string, error)
}
