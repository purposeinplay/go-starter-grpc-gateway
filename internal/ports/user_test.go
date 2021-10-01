package ports

import (
	"context"
	"github.com/pborman/uuid"
	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (ts *StarterTestSuite) TestStarterAPI_FindUserWithUnauthenticated() {
	ctx := context.Background()

	authToken := "123"
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+authToken)

	req := &apigrpc.FindUserRequest{Id: "1"}

	_, err := ts.client.FindUser(ctx, req)
	ts.Error(err)

	st, ok := status.FromError(err)
	ts.True(ok)

	ts.Equal(codes.Unauthenticated, st.Code())
}

func (ts *StarterTestSuite) TestStarterAPI_FindUsers() {
	ctx := context.Background()

	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI5NzU5MDMwMS1iNWJjLTRkM2YtOWRmNS0zNjdhOWMzYjVjMmQiLCJuYW1lIjoiSm9obiBEb2UiLCJpYXQiOjE1MTYyMzkwMjJ9.KAEUKhHbWegoQMA3HHqBvP7KZ3oXn7wdaBDr42PQJ4U"
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+authToken)

	user1 := user.User{
		Base: domain.Base{
			ID:             uuid.Parse("97590301-b5bc-4d3f-9df5-367a9c3b5c2d"),
		},
		Email:             "me@email.com",
	}

	ts.db.Create(&user1)

	req := new(emptypb.Empty)

	res, err := ts.client.FindUsers(ctx, req)
	ts.NoError(err)

	ts.Assert().Equal(1, len(res.GetUsers()))

	u := res.GetUsers()[0]

	ts.Equal("97590301-b5bc-4d3f-9df5-367a9c3b5c2d", u.Id)
	ts.Equal("me@email.com", u.Email)
}

func (ts *StarterTestSuite) TestStarterAPI_CreateUser() {
	ctx := context.Background()

	req := apigrpc.CreateUserRequest{
		Email: "me@email.com",
	}

	res, err := ts.client.CreateUser(ctx, &req)
	ts.NoError(err)

	ts.Assert().Equal("me@email.com", res.GetUser().Email)
}

func (ts *StarterTestSuite) TestStarterAPI_CreateUserWithError() {
	ctx := context.Background()

	req := apigrpc.CreateUserRequest{
		Email: "",
	}

	_, err := ts.client.CreateUser(ctx, &req)
	ts.Error(err)

	st, ok := status.FromError(err)
	ts.True(ok)

	details := st.Details()
	detailsErr := details[0].(*apigrpc.StarterErrorResponse)
	ts.Assert().Equal(apigrpc.ErrorCode_EMAIL_REQUIRED_ERROR, detailsErr.GetError().GetCode())
}
