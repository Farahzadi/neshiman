package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type userCtxKey string

const UserIDKey userCtxKey = "user_id"
const UserRoleKey userCtxKey = "user_role"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			var userID uuid.UUID
			var userRole string
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				// stub: parse JWT and extract user info
				// Phase 6 will implement real JWT validation here
			}
			// stub: default user (uuid.Nil) until real auth is implemented
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(UserIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
