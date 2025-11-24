package blog

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateBlog(name string, ctx context.Context) error {
	var id int

	query := `INSERT INTO blogs (blog_name) VALUES ($1) RETURNING id`

	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return err
	}
	log.Println(id)

	return nil
}

func (r *Repository) UpdateBlog() {

}

func (r *Repository) DeleteBlog() {

}

func (r *Repository) GetBlog() {

}

func (r *Repository) GetBlogs() {

}
