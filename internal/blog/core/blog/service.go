package blog

import (
	"context"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateBlog(ctx context.Context, title string) (int64, error) {
	retId, err := s.repository.CreateBlog(ctx, title)
	if err != nil {
		return retId, fmt.Errorf("failed to create blog: %w", err)
	}
	return retId, nil
}

func (s *Service) UpdateBlog(ctx context.Context, title string, id int64) (int64, error) {
	retId, err := s.repository.UpdateBlog(ctx, title, id)
	if err != nil {
		return retId, fmt.Errorf("failed to update blog: %w", err)
	}
	return retId, nil
}

func (s *Service) DeleteBlog(ctx context.Context, id int64) (int64, error) {
	retId, err := s.repository.DeleteBlog(ctx, id)
	if err != nil {
		return retId, fmt.Errorf("failed to delete blog: %w", err)
	}
	return retId, nil
}

func (s *Service) GetBlog(ctx context.Context, id int64) (Blog, error) {
	blog, err := s.repository.GetBlog(ctx, id)
	if err != nil {
		return blog, fmt.Errorf("failed to get blog: %w", err)
	}
	return blog, nil
}

func (s *Service) GetBlogs(ctx context.Context) ([]Blog, error) {
	blog, err := s.repository.GetBlogs(ctx)
	if err != nil {
		return blog, fmt.Errorf("failed to create blogs: %w", err)
	}
	return blog, nil
}
