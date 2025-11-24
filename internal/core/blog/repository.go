package blog

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

type Blog struct {
	Id   int64  `db:"id"`
	Name string `db:"blog_name"`
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateBlog(ctx context.Context, name string) error {
	var id int

	query := `INSERT INTO blogs (blog_name) VALUES ($1) RETURNING id`

	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return err
	}
	log.Println(id)

	return nil
}

func (r *Repository) UpdateBlog(ctx context.Context, name string, id int) error {
	var returningId int
	query := `UPDATE blogs SET blog_name = $1 WHERE id = $2 RETURNING id`

	err := r.pool.QueryRow(ctx, query, name, id).Scan(&returningId)
	if err != nil {
		return err
	}

	log.Println(returningId)

	return nil

}

func (r *Repository) DeleteBlog(ctx context.Context, id int) error {
	var returningId int

	query := `DELETE FROM blogs WHERE id = $1 RETURNING id`

	err := r.pool.QueryRow(ctx, query, id).Scan(&returningId)
	if err != nil {
		return err
	}
	log.Println(returningId)

	return nil

}

func (r *Repository) GetBlog(ctx context.Context, id int) error {
	var returningId int

	query := `SELECT * FROM blogs WHERE id = $1 RETURNING id`

	err := r.pool.QueryRow(ctx, query, id).Scan(&returningId)
	if err != nil {
		return err
	}
	log.Println(returningId)

	return nil

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
