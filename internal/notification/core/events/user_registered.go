package events

type UserRegistered struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Time   int64  `json:"time"`
}
