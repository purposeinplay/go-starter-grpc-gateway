package grpc

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc/v1"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/auth"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

// FindUsers queries the system for users.
func (s *Server) FindUsers(
	ctx context.Context,
	_ *emptypb.Empty,
) (*startergrpc.FindUsersResponse, error) {
	userID, err := auth.UUIDFromContextJWT(ctx, s.jwtManager)
	if err != nil {
		return nil, fmt.Errorf(
			"user id from context: %w",
			err,
		)
	}

	users, err := s.app.Queries.FindUsers.Handle(
		ctx,
		user.Filter{
			ID: userID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"find users query: %w",
			err,
		)
	}

	resUsers := make([]*startergrpc.User, 0, len(users))

	for _, u := range users {
		resUsers = append(resUsers, &startergrpc.User{
			Id:    u.ID.String(),
			Email: u.Email,
		})
	}

	return &startergrpc.FindUsersResponse{
		Users: resUsers,
	}, nil
}
