package middleware

import (
	"net/http"
)

var roleHierarchy = map[string]int{
	"superadmin": 100,
	"team_admin": 50,
	"viewer":     10,
}

func RequireRole(minimumRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := UserRoleFromContext(r.Context())
			if userRole == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userLevel, ok := roleHierarchy[userRole]
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			requiredLevel, ok := roleHierarchy[minimumRole]
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if userLevel < requiredLevel {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
