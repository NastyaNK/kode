package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"myproject/pkg/config"
	. "myproject/pkg/middleware"
	"myproject/pkg/service"
	"myproject/pkg/utils"
	"net/http"
)

type API interface {
	Authorize(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	AddNote(w http.ResponseWriter, r *http.Request)
	GetAllNotes(w http.ResponseWriter, r *http.Request)
	Run() error
}

type kodeAPI struct {
	config  *config.Api
	handler http.Handler

	middleware Middleware

	authService service.AuthService
	userService service.UserService
	noteService service.NoteService
}

func NewAPI(appConfig *config.Application, userService service.UserService, noteService service.NoteService) API {
	jwt := utils.NewJwtUtils(appConfig)
	auth := service.NewAuthService(userService, jwt)
	middleware := NewMiddleware(appConfig, jwt)

	api := kodeAPI{
		config:      &appConfig.Api,
		userService: userService,
		noteService: noteService,
		authService: auth,
		middleware:  middleware,
	}

	api.handler = NewRouter(&api)
	return &api
}

func (api *kodeAPI) Run() error {
	addr := fmt.Sprintf("%s:%s", api.config.Host, api.config.Port)
	slog.Info("Server started", "address", addr)
	err := http.ListenAndServe(addr, api.handler)
	if err != nil {
		slog.Error("Server stopped unexpectedly", "error", err)
	}
	return err
}

func (api *kodeAPI) responseWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to send JSON response", "error", err, "payload", payload)
	}
}

func (api *kodeAPI) error(w http.ResponseWriter, statusCode int, message string) {
	api.responseWithJson(w, statusCode, map[string]string{"error": message})
}
