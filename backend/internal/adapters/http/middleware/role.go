package middleware

import (
	"context"
	"net/http"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || userRole == "" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			// stub: proper role hierarchy check
			_ = context.WithValue(r.Context(), UserRoleKey, userRole)
			next.ServeHTTP(w, r)
		})
	}
}
