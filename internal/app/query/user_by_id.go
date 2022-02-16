package query

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// UserByIDReadModel represents how the application is querying
// a user.
type UserByIDReadModel interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
}

// UserByIDHandler holds the dependencies for querying a
// user from the system.
type UserByIDHandler struct {
	readModel UserByIDReadModel
}

// MustNewUserByIDHandler returns an initialized UserByIDHandler.
func MustNewUserByIDHandler(
	readModel UserByIDReadModel,
) UserByIDHandler {
	if readModel == nil {
		panic(errors.NewInvalidError("nil read model"))
	}

	return UserByIDHandler{
		readModel: readModel,
	}
}

// Handle queries a user from the system by
// the given ID.
func (s UserByIDHandler) Handle(
	ctx context.Context,
	id uuid.UUID,
) (User, error) {
	u, err := s.readModel.GetUserByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("read model: %w", err)
	}

	return u, nil
}
