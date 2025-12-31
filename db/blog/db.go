package blog

import (
	"context"
	"habr/db/common"
	"habr/internal/blog/config"
)

func Initialize(ctx context.Context, cfg *config.Config) (*common.Database, error) {
	return common.Initialize(ctx, &common.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
}
