package auth

import (
	"context"

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

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userCtxKey).(*User)
	if !ok {
		return nil
	}
	return user
}

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}
