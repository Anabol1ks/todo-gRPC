package db

import (
	"log"

	"github.com/Anabol1ks/todo-gRPC/internal/models"
)

func AutoMigrate() {
	if err := DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatalf("Ошибка при миграции таблиц: %v", err)
	}

	log.Println("Автомиграция таблиц завершена успешно")
}
