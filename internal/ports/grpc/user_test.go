package grpc_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
)

func TestGetUserByIDWithUnauthenticated(t *testing.T) {
	cfg := &config.Config{
		SERVER: struct {
			Port    int    `mapstructure:"port"`
			Address string `mapstructure:"address"`
		}{
			Port:    7350,
			Address: "0.0.0.0",
		},
		JWT: struct {
			Secret          string `mapstructure:"secret"`
			RefreshTokenExp int    `mapstructure:"refresh_token_exp"`
			AccessTokenExp  int    `mapstructure:"access_token_exp"`
		}{
			Secret:          "5649e3d0-7ba4-411d-a721-202c1c626f5c",
			RefreshTokenExp: 3600,
			AccessTokenExp:  900,
		},
	}

	var (
		_, conn = newServerOnRandomPort(t, app.Application{}, cfg)
		client  = apigrpc.NewGoStarterClient(conn)
		ctx     = context.TODO()
		i       = is.New(t)
	)

	_, err := client.GetUser(ctx, &apigrpc.GetUserRequest{
		Id: "1",
	})
	i.NoErr(err)
}

//
// func (ts *StarterTestSuite) TestStarterAPI_FindUserWithUnauthenticated() {
// 	ctx := context.Background()
//
// 	authToken := "123"
// 	ctx = metadata.AppendToOutgoingContext(ctx, "authorization",
// 	"Bearer "+authToken)
//
// 	req := &apigrpc.FindUserRequest{Id: "1"}
//
// 	_, err := ts.client.FindUser(ctx, req)
// 	ts.Error(err)
//
// 	st, ok := status.FromError(err)
// 	ts.True(ok)
//
// 	ts.Equal(codes.Unauthenticated, st.Code())
// }
//
// func (ts *StarterTestSuite) TestStarterAPI_FindUsers() {
// 	ctx := context.Background()
//
// 	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ
// 	zdWIiOiI5NzU5MDMwMS1iNWJjLTRkM2YtOWRmNS0zNjdhOWMzYjVjMmQiLCJuY
// 	W1lIjoiSm9obiBEb2UiLCJpYXQiOjE1MTYyMzkwMjJ9.KAEUKhHbWegoQMA3HH
// 	qBvP7KZ3oXn7wdaBDr42PQJ4U"
// 	ctx = metadata.AppendToOutgoingContext(ctx,
// 	"authorization", "Bearer "+authToken)
//
// 	user1 := user.User{
// 		Base: domain.Base{
// 			ID: uuid.Parse("97590301-b5bc-4d3f-9df5-367a9c3b5c2d"),
// 		},
// 		Email: "me@email.com",
// 	}
//
// 	ts.db.Create(&user1)
//
// 	req := new(emptypb.Empty)
//
// 	res, err := ts.client.FindUsers(ctx, req)
// 	ts.NoError(err)
//
// 	ts.Assert().Equal(1, len(res.GetUsers()))
//
// 	u := res.GetUsers()[0]
//
// 	ts.Equal("97590301-b5bc-4d3f-9df5-367a9c3b5c2d", u.Id)
// 	ts.Equal("me@email.com", u.Email)
// }
//
// func (ts *StarterTestSuite) TestStarterAPI_CreateUser() {
// 	ctx := context.Background()
//
// 	req := apigrpc.CreateUserRequest{
// 		Email: "me@email.com",
// 	}
//
// 	res, err := ts.client.CreateUser(ctx, &req)
// 	ts.NoError(err)
//
// 	ts.Assert().Equal("me@email.com", res.GetUser().Email)
// }
//
// func (ts *StarterTestSuite) TestStarterAPI_CreateUserWithError() {
// 	ctx := context.Background()
//
// 	req := apigrpc.CreateUserRequest{
// 		Email: "",
// 	}
//
// 	_, err := ts.client.CreateUser(ctx, &req)
// 	ts.Error(err)
//
// 	st, ok := status.FromError(err)
// 	ts.True(ok)
//
// 	details := st.Details()
// 	detailsErr := details[0].(*apigrpc.StarterErrorResponse)
// 	ts.Assert().Equal(apigrpc.ErrorCode_EMAIL_REQUIRED_ERROR,
// 		detailsErr.GetError().GetCode())
// }
