package service

import (
	"context"
	"errors"
	"log/slog"
	. "myproject/pkg/models"
	"myproject/pkg/repo"
)

type UserServiceImpl struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id int64) (*User, error) {
	slog.Info("Fetching user by ID", slog.Int64("user_id", id))
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch user", slog.Int64("user_id", id), slog.Any("error", err))
		return nil, err
	}
	slog.Info("User fetched successfully", slog.Int64("user_id", user.Id), slog.String("username", user.Name))
	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *User) error {
	slog.Info("Creating user", slog.String("username", user.Name))
	err := s.repo.InsertUser(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", slog.String("username", user.Name), slog.Any("error", err))
		return err
	}
	slog.Info("User created successfully", slog.Int64("user_id", user.Id), slog.String("username", user.Name))
	return nil
}

func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	slog.Info("Fetching user by username", slog.String("username", username))
	user, err := s.repo.GetUserByName(ctx, username)
	if err != nil {
		slog.Error("Failed to fetch user by username", slog.String("username", username), slog.Any("error", err))
		return nil, err
	}

	// Если пользователь не найден, возвращаем ошибку
	if user == nil || user.Id == 0 {
		err = errors.New("user not found")
		slog.Warn("User not found", slog.String("username", username))
		return nil, err
	}

	slog.Info("User fetched successfully", slog.Int64("user_id", user.Id), slog.String("username", user.Name))
	return user, nil
}
