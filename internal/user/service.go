package user

import (
	"context"
	"errors"
	"log"

	todov1 "github.com/Anabol1ks/todo-gRPC/gen/go/proto/todo"
	"github.com/Anabol1ks/todo-gRPC/internal/auth"
	"github.com/Anabol1ks/todo-gRPC/internal/models"
	"github.com/cockroachdb/errors/grpc/status"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type Service struct {
	todov1.UnimplementedUserServiceServer
	DB  *gorm.DB
	JWT *auth.JWTManager
}

func (s *Service) Register(ctx context.Context, req *todov1.RegisterRequest) (*todov1.AuthResponse, error) {
	op := "Register"
	log.Printf("[%s] start", op)
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	var exists models.User
	if err := s.DB.Where("email = ?", req.Email).First(&exists).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "email already registered")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user := models.User{
		Email:    req.Email,
		Password: string(passHash),
		Nickname: req.Nickname,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	accessToken, refreshToken, err := s.JWT.Generate(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &todov1.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *todov1.LoginRequest) (*todov1.AuthResponse, error) {
	op := "Login"
	log.Printf("[%s] start", op)
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, refreshToken, err := s.JWT.Generate(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &todov1.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, req *todov1.GetProfileRequest) (*todov1.UserResponse, error) {
	op := "GetProfile"
	log.Printf("[%s] start", op)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, status.Error(codes.Internal, "user_id not found in context")
	}
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "db error")
	}

	return &todov1.UserResponse{
		Nickname: user.Nickname,
		Email:    user.Email,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *todov1.RefreshTokenRequest) (*todov1.AuthResponse, error) {
	op := "RefreshToken"
	log.Printf("[%s] start", op)

	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	log.Println(req)

	userID, err := s.JWT.VerifyRefresh(req.RefreshToken)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	newAccess, newRefresh, err := s.JWT.Generate(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate new tokens")
	}

	return &todov1.AuthResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
