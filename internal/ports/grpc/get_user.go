package grpc

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/auth"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// GetUser returns a user from the system.
func (s *Server) GetUser(
	ctx context.Context,
	req *apigrpc.GetUserRequest,
) (*apigrpc.GetUserResponse, error) {
	u, err := auth.UUIDFromContextJWT(ctx, s.jwtManager)
	if err != nil {
		return nil, s.handleErr(fmt.Errorf(
			"authenticate user: %w",
			err,
		))
	}

	if u.String() != req.Id {
		return nil, s.handleErr(errors.NewUnauthorizedError("user"))
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, s.handleErr(fmt.Errorf(
			"parse user id: %w",
			err,
		))
	}

	user, err := s.app.Queries.UserByID.Handle(ctx, userID)
	if err != nil {
		return nil, s.handleErr(fmt.Errorf(
			"user by id query: %w",
			err,
		))
	}

	return &apigrpc.GetUserResponse{
		User: &apigrpc.User{
			Id:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}
