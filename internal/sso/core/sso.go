package core

import (
	"context"
	"errors"
	"log/slog"
	pb_sso "sso/pkg/sso"
)

type AuthService interface {
	Register(ctx context.Context, email string, password string) (userID uint64, err error)
	Login(ctx context.Context, email string, password string, appId uint32) (token string, err error)
}

type Service struct {
	authService AuthService
	pb_sso.UnimplementedSSOServer
	logger *slog.Logger
}

func NewService(authService AuthService, log *slog.Logger) *Service {
	return &Service{
		authService: authService,
		logger:      log,
	}
}

func (s *Service) Register(ctx context.Context, req *pb_sso.RegisterRequest) (*pb_sso.RegisterResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	userID, err := s.authService.Register(ctx, req.Email, req.Password)

	return &pb_sso.RegisterResponse{UserId: userID}, err
}

func (s *Service) Login(ctx context.Context, req *pb_sso.LoginRequest) (*pb_sso.LoginResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}
	if req.AppId == 0 {
		return nil, errors.New("appId is required")
	}

	token, err := s.authService.Login(ctx, req.Email, req.Password, req.AppId)

	return &pb_sso.LoginResponse{Token: token}, err
}
