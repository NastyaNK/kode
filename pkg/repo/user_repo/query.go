package user_repo

import (
	"context"
	"fmt"
	"log/slog"
	. "myproject/pkg/models"
)

func (repo *UserRepository) GetUserById(ctx context.Context, id int64) (*User, error) {
	var user User
	slog.Debug("Fetching user by ID", slog.Int64("user_id", id))

	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		slog.Error("Failed to fetch user by ID", slog.String("error", err.Error()), slog.Int64("user_id", id))
		return &user, fmt.Errorf("не удалось получить пользователя: %w", err)
	}

	slog.Info("User fetched successfully", slog.Int64("user_id", user.Id))
	return &user, nil
}

func (repo *UserRepository) InsertUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`
	slog.Debug("Executing query", slog.String("query", query), slog.String("username", user.Name))

	err := repo.db.QueryRowContext(ctx, query, user.Name, user.Password).Scan(&user.Id)
	if err != nil {
		slog.Error("Failed to insert user", slog.String("error", err.Error()), slog.Any("user", user))
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
	}

	slog.Info("User inserted successfully", slog.Int64("user_id", user.Id))
	return nil
}

func (repo *UserRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	var user User
	slog.Debug("Fetching user by name", slog.String("username", name))

	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE name=$1", name)
	if err != nil {
		slog.Error("Failed to fetch user by name", slog.String("error", err.Error()), slog.String("username", name))
		return &user, fmt.Errorf("не удалось найти пользователя по имени: %w", err)
	}

	slog.Info("User fetched successfully", slog.String("username", user.Name))
	return &user, nil
}
