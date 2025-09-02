package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
	// other relevant user fields
}

func NewUser(id uuid.UUID) *User {
	return &User{
		ID: id,
	}
}

const userCtxKey = "user"

func UserFromContext(c *fiber.Ctx) *User {
	user, ok := c.Locals(userCtxKey).(*User)
	if !ok {
		return nil
	}
	return user
}

func setUser(c *fiber.Ctx, user *User) {
	c.Locals(userCtxKey, user)
}
