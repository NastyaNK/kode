package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	. "myproject/pkg/context"
	. "myproject/pkg/models"
	"myproject/pkg/repo"
	"myproject/pkg/utils"
)

type NoteServiceImpl struct {
	repo repo.NoteRepository
}

func NewNoteService(repo repo.NoteRepository) NoteService {
	return &NoteServiceImpl{
		repo: repo,
	}
}

var (
	ErrNoteService = fmt.Errorf("ошибка сервиса заметок")
)

func (s *NoteServiceImpl) GetAllNotes(ctx context.Context) ([]Note, error) {
	userId := ctx.Value(UserKey).(User).Id
	slog.Debug("Fetching all notes for user", slog.Int64("user_id", userId))

	notes, err := s.repo.GetAllNotes(ctx, userId)
	if err != nil {
		slog.Error("Failed to fetch notes", slog.String("error", err.Error()), slog.Int64("user_id", userId))
		return nil, err
	}

	slog.Info("Notes fetched successfully", slog.Int64("user_id", userId), slog.Int("note_count", len(notes)))
	return notes, nil
}

func (s *NoteServiceImpl) AddNote(ctx context.Context, note *Note) ([]SpellResultTO, error) {
	user := ctx.Value(UserKey).(User)
	note.UserId = user.Id

	slog.Debug("Adding note", slog.Int64("user_id", user.Id), slog.String("content", note.Content))

	checkResult, err := s.Check(note.Content)
	if err != nil || len(checkResult) > 0 {
		if err != nil {
			slog.Error("Spell check failed", slog.String("error", err.Error()))
		} else {
			slog.Warn("Spell check issues detected", slog.Any("check_result", checkResult))
		}
		return utils.ConvertAllSpellResults(checkResult), err
	}

	err = s.repo.InsertNote(ctx, note)
	if err != nil {
		slog.Error("Failed to insert note", slog.String("error", err.Error()), slog.Any("note", note))
		return nil, err
	}

	slog.Info("Note added successfully", slog.Int64("note_id", note.Id))
	return nil, nil
}

func (s *NoteServiceImpl) Check(text string) ([]SpellCheckResult, error) {
	resp, err := http.Post(
		"https://speller.yandex.net/services/spellservice.json/checkText",
		"application/x-www-form-urlencoded",
		strings.NewReader("text="+text),
	)
	if err != nil {
		slog.Error("Spell check request failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrNoteService, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Yandex.Speller responded with non-OK status", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("%w: сервис проверки орфографии временно недоступен", ErrNoteService)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read spell check response body", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrNoteService, err)
	}

	var result []SpellCheckResult
	if err := json.Unmarshal(body, &result); err != nil {
		slog.Error("Failed to unmarshal spell check response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrNoteService, err)
	}

	slog.Debug("Spell check completed successfully", slog.Any("results", result))
	return result, nil
}
