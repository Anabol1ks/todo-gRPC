package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDBPostgres() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s prefer_simple_protocol=true", host, port, user, password, dbname, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных PostgreSQL")
}
