package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
	}
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	// Копируем данные в буфер
	rw.body.Write(p)
	// Отправляем данные клиенту
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)

		slog.Debug("Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(rw, r)

		slog.Debug("Outgoing response",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status", rw.statusCode),
			slog.String("body", rw.body.String()),
		)
	})
}
