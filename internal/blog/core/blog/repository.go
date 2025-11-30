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

func (r *Repository) CreateBlog(ctx context.Context, name string) (int64, error) {
	var id int64

	query := `INSERT INTO blogs (blog_name) VALUES ($1) RETURNING id`

	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return id, err
	}
	log.Println(id)

	return id, nil
}

func (r *Repository) UpdateBlog(ctx context.Context, name string, id int64) (int64, error) {
	var returningId int64
	query := `UPDATE blogs SET blog_name = $1 WHERE id = $2 RETURNING id`

	err := r.pool.QueryRow(ctx, query, name, id).Scan(&returningId)
	if err != nil {
		return returningId, err
	}

	log.Println(returningId)

	return returningId, nil

}

func (r *Repository) DeleteBlog(ctx context.Context, id int64) (int64, error) {
	var returningId int64

	query := `DELETE FROM blogs WHERE id = $1 RETURNING id`

	err := r.pool.QueryRow(ctx, query, id).Scan(&returningId)
	if err != nil {
		return returningId, err
	}
	log.Println(returningId)

	return returningId, nil

}

func (r *Repository) GetBlog(ctx context.Context, id int64) (Blog, error) {
	var blog Blog

	query := `SELECT * FROM blogs WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(&blog.Id, &blog.Name)
	if err != nil {
		return blog, err
	}
	log.Println("выполнен запрос GetBlog", blog.Id, blog.Name)

	return blog, nil

}

func (r *Repository) GetBlogs(ctx context.Context) ([]Blog, error) {
	var blogs []Blog

	query := `SELECT * FROM blogs`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := Blog{}
		err = rows.Scan(&row.Id, &row.Name)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, row)
	}

	return blogs, nil
}
