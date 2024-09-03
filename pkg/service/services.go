package service

import (
	"context"
	. "myproject/pkg/models"
)

type AuthService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, username, password string) (string, error)
	ValidateUser(ctx context.Context, tokenString string) (*User, error)
}

type UserService interface {
	GetUser(context.Context, int64) (*User, error)
	GetUserByUsername(context.Context, string) (*User, error)
	CreateUser(context.Context, *User) error
}

type NoteService interface {
	GetAllNotes(ctx context.Context) ([]Note, error)
	AddNote(ctx context.Context, note *Note) ([]SpellResultTO, error)
}
