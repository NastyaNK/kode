package note_repo

import (
	"context"
	"fmt"
	"log/slog"
	. "myproject/pkg/models"
)

func (repo *NoteRepository) InsertNote(ctx context.Context, note *Note) error {
	query := `INSERT INTO notes (user_id, content) VALUES ($1, $2) RETURNING id`
	slog.Debug("Executing query", slog.String("query", query), slog.Int64("user_id", note.UserId), slog.String("content", note.Content))

	err := repo.db.QueryRowContext(ctx, query, note.UserId, note.Content).Scan(&note.Id)
	if err != nil {
		slog.Error("Failed to insert note", slog.String("error", err.Error()), slog.Any("note", note))
		return fmt.Errorf("не удалось добавить заметку: %w", err)
	}

	slog.Info("Note inserted successfully", slog.Int64("note_id", note.Id))
	return nil
}

func (repo *NoteRepository) GetAllNotes(ctx context.Context, userId int64) ([]Note, error) {
	var notes []Note
	query := `SELECT id, user_id, content FROM notes WHERE user_id = $1`
	slog.Debug("Executing query", slog.String("query", query), slog.Int64("user_id", userId))

	err := repo.db.SelectContext(ctx, &notes, query, userId)
	if err != nil {
		slog.Error("Failed to fetch notes", slog.String("error", err.Error()), slog.Int64("user_id", userId))
		return nil, fmt.Errorf("не удалось получить заметки: %w", err)
	}

	slog.Info("Notes fetched successfully", slog.Int64("user_id", userId), slog.Int("note_count", len(notes)))
	return notes, nil
}
