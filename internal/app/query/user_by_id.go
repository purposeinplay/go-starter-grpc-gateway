package query

import (
	"context"
	"errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain"

	"github.com/pborman/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
)

type UserByIdCmd struct {
	Id string
}

type UserByIdHandler struct {
	logger *zap.Logger
	repo   user.UserRepository
}

func NewUserById(
	logger     *zap.Logger,
	repo user.UserRepository,
) UserByIdHandler {
	return UserByIdHandler{
		logger: logger,
		repo: repo,
	}
}

func (s *UserByIdHandler) First(ctx context.Context, q UserByIdCmd) (*user.User, error) {
	t, err := s.repo.First(ctx, user.User{
		Base:              domain.Base{
			ID: uuid.Parse(q.Id),
		},
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return t, nil
}
