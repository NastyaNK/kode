package repo

import (
	"context"
	"myproject/pkg/config"
	. "myproject/pkg/models"
	"myproject/pkg/repo/connection"
	note "myproject/pkg/repo/note_repo"
	user "myproject/pkg/repo/user_repo"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
}

type NoteRepository interface {
	InsertNote(ctx context.Context, note *Note) error
	GetAllNotes(ctx context.Context, userId int64) ([]Note, error)
}

func NewConnection(config *config.Database) (*connection.Connection, error) {
	return connection.NewConnection(config)
}

func NewUserRepository(connection *connection.Connection) UserRepository {
	return user.NewUserRepository(connection)
}

func NewNoteRepository(connection *connection.Connection) NoteRepository {
	return note.NewNoteRepository(connection)
}
