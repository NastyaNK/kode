package user_repo

import (
	"myproject/pkg/repo/connection"
)

type UserRepository struct {
	db *connection.Connection
}

func NewUserRepository(connection *connection.Connection) *UserRepository {
	return &UserRepository{connection}
}
