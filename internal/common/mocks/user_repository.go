package mocks

import (
	"context"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

var (
	_ user.Repository          = (*UserRepository)(nil)
	_ query.FindUsersReadModel = (*UserRepository)(nil)
	_ query.UserByIDReadModel  = (*UserRepository)(nil)
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)

	return args.Error(0)
}

func (m *UserRepository) FindUsers(
	ctx context.Context,
	filter user.Filter,
) ([]query.User, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).([]query.User), args.Error(1)
}

func (m *UserRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (query.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(query.User), args.Error(1)
}
