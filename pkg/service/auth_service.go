package service

import (
	"context"
	"errors"
	"log/slog"
	. "myproject/pkg/models"
	"myproject/pkg/utils"
)

type AuthServiceImpl struct {
	userService UserService
	jwtService  *utils.JwtUtils
}

func NewAuthService(userService UserService, jwtService *utils.JwtUtils) AuthService {
	return &AuthServiceImpl{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, user *User) error {
	slog.Info("Attempting to register user",
		"username", user.Name,
	)
	existingUser, err := s.userService.GetUserByUsername(ctx, user.Name)
	if err == nil && existingUser.Id != 0 {
		slog.Warn("User already exists",
			"username", user.Name,
		)
		return errors.New("user already exists")
	}
	if err := s.userService.CreateUser(ctx, user); err != nil {
		slog.Error("Failed to register user",
			"error", err,
			"username", user.Name,
		)
		return err
	}
	slog.Info("User registered successfully",
		"user_id", user.Id,
		"username", user.Name,
	)
	return nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (string, error) {
	slog.Info("Attempting to log in user",
		"username", username,
	)
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil || user.Password != password {
		slog.Warn("Invalid username or password",
			"username", username,
		)
		return "", errors.New("invalid username or password")
	}

	claims := map[string]interface{}{
		"id":   user.Id,
		"name": user.Name,
	}

	token, err := s.jwtService.GenerateToken(claims)
	if err != nil {
		slog.Error("Failed to generate token",
			"error", err,
			"username", username,
		)
		return "", err
	}

	slog.Info("User logged in successfully",
		"username", username,
		"token", token,
	)
	return token, nil
}

func (s *AuthServiceImpl) ValidateUser(ctx context.Context, tokenString string) (*User, error) {
	slog.Info("Validating user token")
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		slog.Error("Token validation failed",
			"error", err,
		)
		return nil, err
	}

	userId, ok := claims["id"].(float64)
	if !ok {
		slog.Error("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	user, err := s.userService.GetUser(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to retrieve user from token",
			"error", err,
			"user_id", userId,
		)
		return nil, err
	}

	slog.Info("User validated successfully",
		"user_id", user.Id,
		"username", user.Name,
	)
	return user, nil
}
