package blog

type Blog struct {
	Id   int64  `db:"id"`
	Name string `db:"blog_name"`
}
