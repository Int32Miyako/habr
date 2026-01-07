package events

type UserRegistered struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Time   int64  `json:"time"`
}
