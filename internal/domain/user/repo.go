package user

import (
	"context"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// Repository defines methods for User CRUD actions.
type Repository interface {
	CreateUser(ctx context.Context, u *User) error
}

// Filter represents the data that can be used for
// filtering users found in the repository.
type Filter struct {
	ID    uuid.UUID
	Email *string
}
