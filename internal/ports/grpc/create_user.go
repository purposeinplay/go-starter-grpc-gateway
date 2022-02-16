package grpc

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// CreateUser adds a new user to the system.
func (s *Server) CreateUser(
	ctx context.Context,
	req *apigrpc.CreateUserRequest,
) (*apigrpc.CreateUserResponse, error) {
	newUserID := uuid.New()

	err := s.app.Commands.CreateUser.Handle(ctx, command.CreateUser{
		ID:    newUserID,
		Email: req.Email,
	})
	if err != nil {
		return nil, s.handleErr(fmt.Errorf(
			"create user command: %w",
			err,
		))
	}

	return &apigrpc.CreateUserResponse{
		User: &apigrpc.User{
			Id:    newUserID.String(),
			Email: req.Email,
		},
	}, nil
}
