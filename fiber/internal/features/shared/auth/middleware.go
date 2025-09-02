package auth

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Authenticate(loggerProvider func(name string) *slog.Logger) fiber.Handler {
	logger := loggerProvider("authenticate_middleware")
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")

		if strings.HasPrefix(header, "Bearer ") {
			token := strings.TrimPrefix(header, "Bearer ")

			// fake token with id included
			if strings.HasPrefix(token, "TEST_TOKEN_") {
				idString := strings.TrimPrefix(token, "TEST_TOKEN_")
				id, err := uuid.Parse(idString)
				if err != nil {
					logger.DebugContext(c.UserContext(), "failed to parse user id from token", slog.Any("error", err))
				} else {
					user := NewUser(id)
					setUser(c, user)
					logger.DebugContext(c.UserContext(), "authenticated user", slog.String("user_id", idString))
				}
			}
		}

		return c.Next()
	}
}

func RequireAuthentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := UserFromContext(c)
		if user == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func RequirePermission(permission string, loggerProvider func(name string) *slog.Logger) fiber.Handler {
	logger := loggerProvider("require_permission_middleware")
	return func(c *fiber.Ctx) error {
		// TODO: check if the user has the required permission
		logger.DebugContext(c.UserContext(), "Checking permission", slog.String("permission", permission))
		return c.Next()
	}
}
