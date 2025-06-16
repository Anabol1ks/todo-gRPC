package main

import (
	"todo-grpc/internal/config"
	"todo-grpc/internal/db"
)

func main() {
	config.InitConfig()
	db.ConnectDBPostgres()
	db.AutoMigrate()
}
