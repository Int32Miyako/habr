package dto

// Auth DTOs
type (
	RequestRegisterUser struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ResponseRegisterUser struct {
		UserId int64 `json:"user_id"`
	}

	RequestLoginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginUserDto struct {
		AccessToken  string
		RefreshToken string
		UserId       int64
	}

	ResponseLoginUser struct {
		AccessToken string `json:"access_token"`
		UserId      int64  `json:"user_id"`
	}
)
