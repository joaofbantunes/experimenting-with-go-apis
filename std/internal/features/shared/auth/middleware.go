package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

func RequireAuthentication() shared.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}

			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")
				if token != "TEST_TOKEN" {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				}
			} else {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermission(permission string, loggerProvider func(name string) *slog.Logger) shared.Middleware {
	logger := loggerProvider("require_permission_middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check if the user has the required permission
			logger.Debug("Checking permission", slog.String("permission", permission))
			next.ServeHTTP(w, r)
		})
	}
}
