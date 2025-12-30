package main

import (
	"fmt"
	"habr/internal/notification/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.GRPC)
}
