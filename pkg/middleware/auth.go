package middleware

import (
	"context"
	"log/slog"
	. "myproject/pkg/context"
	. "myproject/pkg/models"
	"myproject/pkg/utils"
	"net/http"
	"strings"
)

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Warn("Authorization token missing",
				"endpoint", r.URL.Path,
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Authorization token missing", http.StatusForbidden)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			slog.Warn("Invalid Authorization header format",
				"endpoint", r.URL.Path,
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Invalid Authorization header format", http.StatusForbidden)
			return
		}
		token := tokenParts[1]

		data, err := m.jwtUtils.ValidateToken(token)
		if err != nil {
			slog.Error("Token validation failed",
				"error", err,
				"token", token,
				"endpoint", r.URL.Path,
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Token validation failed", http.StatusForbidden)
			return
		}

		var user User
		err = utils.ReMarshal(data, &user)
		if err != nil {
			slog.Error("Failed to unmarshal token data into user struct",
				"error", err,
				"data", data,
				"endpoint", r.URL.Path,
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Failed to extract user data from token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		slog.Info("User authenticated",
			"user", user.Name,
			"endpoint", r.URL.Path,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
