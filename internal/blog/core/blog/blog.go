package blog

type Blog struct {
	Id   int64  `json:"id" dbcodes:"id"`
	Name string `json:"blog_name" dbcodes:"blog_name"`
}
