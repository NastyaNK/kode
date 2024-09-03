package api

import (
	"encoding/json"
	"log/slog"
	. "myproject/pkg/models"
	"net/http"
)

func (api *kodeAPI) AddNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		slog.Error("Failed to decode note data", "error", err, "endpoint", r.URL.Path)
		api.error(w, http.StatusBadRequest, "Invalid request data format")
		return
	}

	results, err := api.noteService.AddNote(r.Context(), &note)

	status := http.StatusOK
	response := map[string]interface{}{}

	existErrorsInText := len(results) > 0
	if err != nil || existErrorsInText {
		status = http.StatusBadRequest
		if existErrorsInText {
			response["error"] = "Spelling errors detected! Please correct them to save the note!"
			response["spell"] = results
			slog.Warn("Note contains spelling errors", "note", note.Content, "endpoint", r.URL.Path)
		} else {
			response["error"] = err.Error()
		}
		slog.Error("Failed to add note", "note", note.Content, "error", err)
	} else {
		response["message"] = "Note added successfully"
		response["note"] = note
		slog.Info("Note added successfully", "note", note.Content, "endpoint", r.URL.Path)
	}

	api.responseWithJson(w, status, response)
}

func (api *kodeAPI) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := api.noteService.GetAllNotes(r.Context())
	if err != nil {
		slog.Error("Failed to retrieve notes", "error", err, "endpoint", r.URL.Path)
		api.error(w, http.StatusInternalServerError, "Failed to retrieve notes: "+err.Error())
		return
	}

	slog.Info("Notes retrieved successfully", "count", len(notes), "endpoint", r.URL.Path)
	api.responseWithJson(w, http.StatusOK, notes)
}
