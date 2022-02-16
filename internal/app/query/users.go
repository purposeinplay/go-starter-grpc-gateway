package query

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
)

// FindUsersReadModel represents the application can query
// for users.
type FindUsersReadModel interface {
	FindUsers(
		ctx context.Context,
		filter user.Filter,
	) ([]User, error)
}

// FindUsersHandler holds the dependencies for querying
// users from the system.
type FindUsersHandler struct {
	readModel FindUsersReadModel
}

// MustNewFindUsersHandler returns an initialized FindUsersHandler.
func MustNewFindUsersHandler(
	readModel FindUsersReadModel,
) FindUsersHandler {
	if readModel == nil {
		panic(errors.NewInvalidError("nil read model"))
	}

	return FindUsersHandler{
		readModel: readModel,
	}
}

// Handle queries the system for the given
// users based on the filter provided.
func (s FindUsersHandler) Handle(
	ctx context.Context,
	filter user.Filter,
) ([]User, error) {
	users, err := s.readModel.FindUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("read model: %w", err)
	}

	return users, nil
}
