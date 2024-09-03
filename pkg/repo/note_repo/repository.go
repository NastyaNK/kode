package note_repo

import (
	"myproject/pkg/repo/connection"
)

type NoteRepository struct {
	db *connection.Connection
}

func NewNoteRepository(connection *connection.Connection) *NoteRepository {
	return &NoteRepository{connection}
}
