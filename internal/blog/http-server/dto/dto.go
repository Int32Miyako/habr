package dto

type (
	RequestCreateBlog struct {
		Name string `json:"name"`
	}
	ResponseCreateBlog struct {
		Id int64 `json:"id"`
	}

	RequestUpdateBlog struct {
		Name string `json:"name"`
	}

	ResponseUpdateBlog struct {
		Id int64 `json:"id"`
	}

	ResponseDeleteBlog struct {
		Id int64 `json:"id"`
	}
)

// Auth DTOs
type (
	RequestRegisterUser struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	ResponseRegisterUser struct {
		UserId int64 `json:"user_id"`
	}

	RequestLoginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ResponseLoginUser struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserId       int64  `json:"user_id"`
	}
)
