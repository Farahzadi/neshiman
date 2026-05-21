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
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			// stub: parse JWT and extract user info
			ctx := context.WithValue(r.Context(), UserIDKey, uuid.Nil)
			ctx = context.WithValue(ctx, UserRoleKey, "")
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
