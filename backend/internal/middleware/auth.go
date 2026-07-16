package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	apierrors "taskflow-backend/internal/errors"
	"taskflow-backend/internal/models"
)

type contextKey string

const userContextKey contextKey = "currentUser"

// RequireAuth validates the JWT Bearer token, loads the corresponding user,
// and stores it in the request context. Missing/invalid/expired tokens all
// result in 401 TOKEN_EXPIRED, matching the Python implementation's behavior.
func RequireAuth(db *sqlx.DB, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				WriteError(w, apierrors.TokenExpired())
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, ok := DecodeAccessToken(token, jwtSecret)
			if !ok {
				WriteError(w, apierrors.TokenExpired())
				return
			}

			var user models.User
			err := db.Get(&user, db.Rebind("SELECT * FROM users WHERE id = ?"), userID)
			if err == sql.ErrNoRows {
				WriteError(w, apierrors.TokenExpired())
				return
			}
			if err != nil {
				WriteError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return user
}

func CORS(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, o := range origins {
		allowed[o] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
