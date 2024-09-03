package api

import (
	"encoding/json"
	"log/slog"
	. "myproject/pkg/models"
	"net/http"
)

func (api *kodeAPI) Authorize(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		slog.Warn("Missing authorization data", "endpoint", r.URL.Path, "method", r.Method)
		api.error(w, http.StatusUnauthorized, "Missing authorization data")
		return
	}

	token, err := api.authService.Login(r.Context(), username, password)
	if err != nil {
		slog.Warn("Failed login attempt", "username", username, "error", err)
		api.error(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	api.responseWithJson(w, http.StatusOK, map[string]string{"token": token})
	slog.Info("User authorized", "username", username, "endpoint", r.URL.Path)
}

func (api *kodeAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("Failed to decode user data", "error", err, "endpoint", r.URL.Path)
		api.error(w, http.StatusBadRequest, "Invalid request data format")
		return
	}

	err := api.authService.Register(r.Context(), &user)
	if err != nil {
		slog.Error("User registration failed", "username", user.Name, "error", err)
		api.error(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	api.responseWithJson(w, http.StatusOK, map[string]string{"message": "User " + user.Name + " created successfully"})
	slog.Info("User registered successfully", "username", user.Name, "endpoint", r.URL.Path)
}
