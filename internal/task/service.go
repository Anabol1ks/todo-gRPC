package task

import (
	"context"
	"log"
	"time"

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

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

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
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

func (s *Service) GetTasks(ctx context.Context, req *todov1.GetTasksRequest) (*todov1.TasksList, error) {
	op := "GetTasks"
	log.Printf("[%s] start", op)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}

	var tasks []models.Task
	if err := s.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve tasks")
	}

	var taskResponses []*todov1.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, &todov1.TaskResponse{
			Id:          uint64(task.ID),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserId:      uint64(task.UserID),
		})
	}

	return &todov1.TasksList{
		Tasks: taskResponses,
	}, nil
}

func (s *Service) GetTask(ctx context.Context, req *todov1.GetTaskRequest) (*todov1.TaskResponse, error) {
	op := "GetTask"
	log.Printf("[%s] start", op)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}

	var task models.Task
	if err := s.DB.Where("id = ? AND user_id = ?", req.Id, userID).First(&task).Error; err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	return &todov1.TaskResponse{
		Id:          uint64(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserId:      uint64(task.UserID),
	}, nil
}

func (s *Service) DeleteTask(ctx context.Context, req *todov1.DeleteTaskRequest) (*todov1.Empty, error) {
	op := "DeleteTask"
	log.Printf("[%s] start", op)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}

	var task models.Task
	if err := s.DB.Where("id = ? AND user_id = ?", req.Id, userID).First(&task).Error; err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	if err := s.DB.Delete(&task).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	return &todov1.Empty{
		Value: "Task deleted successfully",
	}, nil
}

func (s *Service) UpdateTask(ctx context.Context, req *todov1.UpdateTaskRequest) (*todov1.TaskResponse, error) {
	op := "UpdateTask"
	log.Printf("[%s] start", op)

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}

	var task models.Task
	if err := s.DB.Where("id = ? AND user_id = ?", req.Id, userID).First(&task).Error; err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.UpdatedAt = time.Now()

	if err := s.DB.Save(&task).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	return &todov1.TaskResponse{
		Id:          uint64(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserId:      uint64(task.UserID),
	}, nil
}
