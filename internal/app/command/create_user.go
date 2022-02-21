package command

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
)

// CreateUser represents the data required
// in order to add a new user into the system.
type CreateUser struct {
	ID    uuid.UUID
	Email string
}

// CreateUserHandler holds the dependencies for adding a
// new user to the system.
type CreateUserHandler struct {
	userRepo user.Repository
}

// MustNewCreateUserHandler returns an initialized CreateUserHandler.
func MustNewCreateUserHandler(
	userRepo user.Repository,
) CreateUserHandler {
	if userRepo == nil {
		panic(errors.NewInvalidError("nil user repo"))
	}

	return CreateUserHandler{
		userRepo: userRepo,
	}
}

// Handle executes the CreateUser command.
func (s *CreateUserHandler) Handle(
	ctx context.Context,
	cmd CreateUser,
) error {
	newUser, err := user.New(
		cmd.ID,
		cmd.Email,
	)
	if err != nil {
		return fmt.Errorf("new user: %w", err)
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
