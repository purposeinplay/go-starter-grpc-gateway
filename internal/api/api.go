package api

import (
	"context"
	"fmt"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/repository"

	grpccommons "github.com/purposeinplay/go-commons/grpc"

	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/purposeinplay/go-commons/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"gorm.io/gorm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
)

type API struct {
	starterapi.UnimplementedGoStarterServer
	config     *config.Config
	logger     *zap.Logger
	grpcServer *grpccommons.Server
	app app.Application
}

func NewAPI(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	db *gorm.DB,
	jwtManager *auth.JWTManager,
) *API {

	const servicePath = "/starter.apigrpc.GoStarter/"
	authRoles := map[string][]string{
		servicePath + "FindUsers":  {"user"},
		servicePath + "FindUser": {"user"},
		servicePath + "CreateUser": {"user"},
	}

	authInterceptor := auth.NewAuthInterceptor(logger, jwtManager, authRoles)

	grpc_zap.ReplaceGrpcLoggerV2(logger)

	serverOptions := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			authInterceptor.Unary(),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger),
		),
	}

	repo := repository.NewUserRepository(db)

	api := &API{
		app: app.Application{
			Commands: app.Commands{
				CreateUser: command.NewCreateUserHandler(logger, repo),
			},
			Queries:  app.Queries{
				FindUsers: query.NewFindUsersHandler(logger, repo),
				UserByID: query.NewUserById(logger, repo),
			},
		},
		config:     config,
		logger:     logger,
	}

	options := []grpccommons.ServerOption{
		grpccommons.WithGrpcServerOptions(serverOptions),

		grpccommons.ReplaceLogger(logger),

		grpccommons.RegisterServer(
			func(server *grpc.Server) {
				starterapi.RegisterGoStarterServer(server, api)
			},
		),

		grpccommons.RegisterGateway(
			func(mux *runtime.ServeMux, dialOptions []grpc.DialOption) {
				dialAddr := fmt.Sprintf("127.0.0.1:%d", config.SERVER.Port-1)
				if config.SERVER.Address != "" {
					dialAddr = fmt.Sprintf("%v:%d", config.SERVER.Address, config.SERVER.Port-1)
				}
				err := starterapi.RegisterGoStarterHandlerFromEndpoint(ctx, mux, dialAddr, dialOptions)
				if err != nil {
					logger.Fatal("connecting to gRPC gateway", zap.Error(err))
				}
			},
		),
	}

	grpcServer := grpccommons.NewServer(options...)
	api.grpcServer = grpcServer

	return api
}

func (a *API) Stop() {
	a.grpcServer.Stop()
}

func (a *API) Healthcheck(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
