package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(api *kodeAPI) http.Handler {
	router := chi.NewRouter()

	router.Use(api.middleware.LoggingMiddleware)

	// Не требуется авторизация
	router.Route("/", func(r chi.Router) {
		r.Post("/reg", api.Register)
		r.Post("/auth", api.Authorize)
	})

	// Требуется авторизация
	router.Route("/note", func(r chi.Router) {
		r.Use(api.middleware.AuthMiddleware)
		r.Get("/all", api.GetAllNotes)
		r.Post("/add", api.AddNote)
	})

	return router
}
