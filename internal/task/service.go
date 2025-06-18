package task

import (
	"context"
	"log"

	todov1 "github.com/Anabol1ks/todo-gRPC/gen/go/proto/todo"
	"github.com/Anabol1ks/todo-gRPC/internal/auth"
	"github.com/Anabol1ks/todo-gRPC/internal/models"
	"github.com/cockroachdb/errors/grpc/status"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type Service struct {
	todov1.UnimplementedTaskServiceServer
	DB  *gorm.DB
	JWT *auth.JWTManager
}

func (s *Service) CreateTask(ctx context.Context, req *todov1.CreateTaskRequest) (*todov1.TaskResponse, error) {
	op := "CreateTask"
	log.Printf("[%s] start", op)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve user")
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		UserID:      user.ID,
	}

	if err := s.DB.Create(&task).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	log.Println(task)
	return &todov1.TaskResponse{
		Id:          uint64(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserId:      uint64(task.UserID),
	}, nil

}

// func (s *Service) GetTask(ctx context.Context, req *todov1.GetTasksRequest) (*todov1.TaskResponse, error) {
// 	op := "GetTask"
// 	log.Printf("[%s] start", op)

// 	userID, ok := ctx.Value("user_id").(uint64)
// 	if !ok {
// 		return nil, status.Error(codes.Internal, "user_id not found in context")
// 	}

// }
