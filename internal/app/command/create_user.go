package command

import (
	"context"

	"go.uber.org/zap"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
)

type CreateUserCmd struct {
	Email          string
}

type CreateUserHandler struct {
	logger            *zap.Logger
	repo              user.UserRepository
}

func NewCreateUserHandler(
	logger     *zap.Logger,
	repo user.UserRepository,
) CreateUserHandler {
	return CreateUserHandler{
		logger: logger,
		repo: repo,
	}
}

func (s *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCmd) (*user.User, error) {

	if cmd.Email == "" {
		return nil, user.ErrEmailIsRequired
	}

	u, err := s.repo.Create(ctx, &user.User{
		Email:             cmd.Email,
	})

	if err != nil {
		return nil, user.ErrCouldNotCreateUser
	}

	return u, nil
}
