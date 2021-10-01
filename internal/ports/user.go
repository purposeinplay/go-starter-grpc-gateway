package ports

import (
	"context"
	"errors"
	"github.com/pborman/uuid"
	"github.com/purposeinplay/go-commons/auth"
	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *grpcServer) FindUsers(ctx context.Context, _ *emptypb.Empty) (*starterapi.FindUsersResponse, error) {

	md, _ := metadata.FromIncomingContext(ctx)

	sign, err := auth.ExtractTokenFromMetadata(md)

	if err != nil {
		a.logger.Error("auth.ExtractTokenFromMetadata error", zap.Error(err))
		return nil, status.Error(codes.Internal, "an error occurred")
	}

	claims, err := a.jwtManager.Verify(sign)

	if err != nil {
		a.logger.Error("jwtManager.Verify error", zap.Error(err))
		return nil, status.Error(codes.Internal, "an error occurred")
	}

	users, err := a.app.Queries.FindUsers.Handle(ctx, query.FindUsersCmd{uuid.Parse(claims.Subject)})

	if err != nil {
		a.logger.Error("queries.FindUsers.Handle error", zap.Error(err))
		st := status.New(codes.Internal, "an error occurred")
		return nil, st.Err()
	}

	res := userEntitiesToResponse(users)
	return res, nil
}

func (a *grpcServer) CreateUser(ctx context.Context, req *starterapi.CreateUserRequest) (*starterapi.CreateUserResponse, error) {
	cmd := command.CreateUserCmd{
		Email:          req.Email,
	}

	u, err := a.app.Commands.CreateUser.Handle(ctx, cmd)

	if err != nil {
		a.logger.Error("tournament.CreateTournament error", zap.Error(err))

		var errCommon *user.Error
		if errors.As(err, &errCommon) {
			st := status.New(codes.InvalidArgument, "could not create user")
			detail := &starterapi.StarterErrorResponse{
				Error: &starterapi.ErrorCode{
					Code: user.GRPCErrorFromError(err),
				},
				Message: err.Error(),
			}

			st, _ = st.WithDetails(detail)
			return nil, st.Err()
		} else {
			return nil, status.Error(codes.Internal, "an error occurred")
		}
	}

	res := &starterapi.CreateUserResponse{User: &starterapi.User{
		Id:    u.ID.String(),
		Email: u.Email,
	}}

	return res, nil
}

func (a *grpcServer) FindUser(ctx context.Context, req *starterapi.FindUserRequest) (*starterapi.FindUserResponse, error) {
	return &starterapi.FindUserResponse{User: &starterapi.User{
		Id:    "123",
		Email: "vlad@asd.com",
	}}, nil
}

func userEntitiesToResponse(users *[]user.User) *starterapi.FindUsersResponse {
	var res starterapi.FindUsersResponse
	for _, u := range *users {
		res.Users = append(res.Users, &starterapi.User{
			Id:    u.ID.String(),
			Email: u.Email,
		})
	}

	return &res
}