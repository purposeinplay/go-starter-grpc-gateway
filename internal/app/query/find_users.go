package query

import (
	"context"

	"go.uber.org/zap"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
)

type FindUsersCmd struct {
}

type FindUsersHandler struct {
	logger *zap.Logger
	repo   user.UserRepository
}

func NewFindUsersHandler(
	logger     *zap.Logger,
	repo user.UserRepository,
) FindUsersHandler {
	return FindUsersHandler{
		logger: logger,
		repo: repo,
	}
}

func (s *FindUsersHandler) Handle(ctx context.Context, q FindUsersCmd) (*[]user.User, error) {
	return s.repo.Find(ctx, user.User{})
}
