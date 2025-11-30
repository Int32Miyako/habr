package blog

type Blog struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"blog_name" db:"blog_name"`
}
