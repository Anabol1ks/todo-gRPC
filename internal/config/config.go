package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}
