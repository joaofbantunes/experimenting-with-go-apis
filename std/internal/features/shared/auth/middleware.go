package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

func Authenticate(loggerProvider func(name string) *slog.Logger) shared.Middleware {
	logger := loggerProvider("authenticate_middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")

				// fake token with id included
				if strings.HasPrefix(token, "TEST_TOKEN_") {
					idString := strings.TrimPrefix(token, "TEST_TOKEN_")
					id, err := uuid.Parse(idString)
					if err != nil {
						logger.DebugContext(r.Context(), "failed to parse user id from token", slog.Any("error", err))
					} else {
						user := NewUser(id)
						r = r.WithContext(ContextWithUser(r.Context(), user))
						logger.DebugContext(r.Context(), "authenticated user", slog.Any("user_id", id))
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthentication() shared.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
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
