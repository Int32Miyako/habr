package notification

import (
	"context"
	"habr/db/common"
	"habr/internal/notification/config"
)

// Initialize инициализирует подключение к базе данных и паникует в случае ошибки
// держит таймаут до 30 секунд для инициализации
func Initialize(ctx context.Context, cfg *config.Config) *common.Database {
	db, err := common.Initialize(ctx, &common.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		panic(err)
	}

	return db
}
