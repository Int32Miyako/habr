package notification

import (
	"context"
	"habr/db/common"
	"habr/internal/notification/config"
	"time"
)

// MustInitialize инициализирует подключение к базе данных и паникует в случае ошибки
// держит таймаут до 30 секунд для инициализации
func MustInitialize(cfg *config.Config) *common.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := common.Initialize(ctx, &common.Config{
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
