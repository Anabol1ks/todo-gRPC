package main

import (
	"log"
	"net"
	"os"

	todov1 "github.com/Anabol1ks/todo-gRPC/gen/go/proto/todo"
	"github.com/Anabol1ks/todo-gRPC/internal/auth"
	"github.com/Anabol1ks/todo-gRPC/internal/config"
	"github.com/Anabol1ks/todo-gRPC/internal/db"
	"github.com/Anabol1ks/todo-gRPC/internal/user"
	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()
	db.ConnectDBPostgres()
	db.AutoMigrate()

	jwtManager := &auth.JWTManager{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		AccessTTL:     config.AccessTTL,
		RefreshTTL:    config.RefreshTTL,
	}

	grpcServer := grpc.NewServer()

	userService := &user.Service{
		DB:  db.DB,
		JWT: jwtManager,
	}

	todov1.RegisterUserServiceServer(grpcServer, userService)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT")) // напр. ":50051"
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server started on %s", os.Getenv("GRPC_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
