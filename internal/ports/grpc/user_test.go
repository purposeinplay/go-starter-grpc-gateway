package grpc_test

import (
	"context"
	"github.com/matryer/is"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/mocks"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/purposeinplay/go-starter-grpc-gateway/apigrpc/v1"
)

func TestServer_CreateUser(t *testing.T) {
	t.Parallel()

	var (
		i   = is.New(t)
		ctx = context.TODO()
	)

	tests := map[string]struct {
		user            *user.User
		createUserErr   error
		expectedGrpcErr error
		reportErr       error
	}{
		"Success": {
			user:            user.MustNew(testID, "user@email.com"),
			createUserErr:   nil,
			expectedGrpcErr: nil,
			reportErr:       nil,
		},
		"Error_InvalidID": {
			user:            user.MustNew(testID, "user@email.com"),
			createUserErr:   errors.NewInvalidError("invalid id"),
			expectedGrpcErr: status.Error(codes.InvalidArgument, "invalid id"),
			reportErr:       nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				mockUserRepo      = new(mocks.UserRepository)
				mockReportService = new(mocks.ReportService)
			)

			_, conn := newTestServer(
				t,
				app.Application{
					Commands: app.Commands{
						CreateUser:  command.MustNewCreateUserHandler(mockUserRepo),
						ReportError: command.MustNewReportErrorHandler(mockReportService),
					},
				},
			)

			client := startergrpc.NewGoStarterClient(conn)

			mockUserRepo.
				On(
					"CreateUser",
					mock.AnythingOfType("*context.valueCtx"),
					test.user,
				).
				Return(
					test.createUserErr,
				)

			mockReportService.
				On(
					"ReportError",
					mock.Anything,
					mock.Anything,
				).
				Return(
					test.reportErr,
				).
				Maybe()

			_, err := client.CreateUser(ctx, &startergrpc.CreateUserRequest{
				Email: test.user.Email(),
			})
			if err != nil {
				st, ok := status.FromError(err)
				i.True(ok)

				expSt, ok := status.FromError(test.expectedGrpcErr)
				i.True(ok)

				i.Equal(expSt.Code(), st.Code())
				i.Equal(expSt.Message(), st.Message())
			}
		})
	}

}
