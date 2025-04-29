package auth

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/sso/auth/models"
	"time"
)

type Repository interface {
	Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Service struct {
	repository Repository
	logger     *slog.Logger
}

func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) Register(ctx context.Context, email string, password string) (userID uint64, err error) {
	const operation = "auth.service.Register"
	logger := s.logger.With("operation", operation)

	logger.Info("Registering user", "email", email, "password", password)
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	if err := s.repository.ExecQueryRow(ctx,
		"insert into users (email, pass_hash) values ($1, $2) returning u_id",
		email, passHash).Scan(&userID); err != nil {
		return 0, err
	}
	logger.Info("User created", "id", userID)

	return userID, nil
}

func (s *Service) Login(ctx context.Context, email string, password string, appId uint32) (token string, err error) {
	const operation = "auth.service.Login"
	logger := s.logger.With("operation", operation)

	logger.Info("Authenticating user", "email", email)

	var user models.User

	if err = s.repository.Get(ctx, &user, "select * from users where email = $1", email); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", err
	}

	var app models.App
	if err := s.repository.Get(ctx, &app, "select * from apps where a_id = $1", appId); err != nil {
		return "", err
	}

	ttl, err := time.ParseDuration("10h")
	if err != nil {
		return "", err
	}
	token, err = NewToken(user, app, ttl)
	if err != nil {
		return "", err
	}

	return token, nil
}
