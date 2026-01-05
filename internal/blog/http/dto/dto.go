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
