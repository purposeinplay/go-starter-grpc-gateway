package api

import (
	"context"
	"fmt"
	grpccommons "github.com/purposeinplay/go-commons/grpc"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/purposeinplay/go-commons/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"gorm.io/gorm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/repository"

	"google.golang.org/grpc"

	"github.com/purposeinplay/go-commons/logs"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/storage/dialer"
)


type StarterTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *grpccommons.Server
	client starterapi.GoStarterClient
}

func (ts *StarterTestSuite) SetupTest() {
	ts.db.Exec("TRUNCATE TABLE users CASCADE")
}

func TestStarterTestSuite(t *testing.T) {
	srv, db := CreateTestAPI(t)
	client := CreateTestAPIClient(t)

	ts := &StarterTestSuite{
		db:     db,
		server: srv,
		client: client,
	}

	suite.Run(t, ts)
}

func CreateTestAPI(t *testing.T) (*grpccommons.Server, *gorm.DB) {
	ctx := context.Background()
	cfg, err := config.LoadTestConfig("../../config.test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	logger := logs.NewLogger()

	db, err := dialer.Connect(cfg)

	if err != nil {
		t.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.AccessTokenExp))

	const servicePath = "/starter.apigrpc.GoStarter/"
	authRoles := map[string][]string{
		servicePath + "FindUser":  {"user"},
	}

	authInterceptor := auth.NewAuthInterceptor(logger, jwtManager, authRoles)

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
		config:     cfg,
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
				dialAddr := fmt.Sprintf("127.0.0.1:%d", cfg.SERVER.Port-1)
				err := starterapi.RegisterGoStarterHandlerFromEndpoint(ctx, mux, dialAddr, dialOptions)
				if err != nil {
					logger.Fatal("connecting to gRPC gateway", zap.Error(err))
				}
			},
		),
	}

	grpcServer := grpccommons.NewServer(options...)
	api.grpcServer = grpcServer

	return grpcServer, db
}

func CreateTestAPIClient(t *testing.T) starterapi.GoStarterClient {
	cfg, err := config.LoadTestConfig("../../config.test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", cfg.SERVER.Address, cfg.SERVER.Port-1),
		grpc.WithInsecure(),
	)

	if err != nil {
		t.Fatal(err)
	}

	return starterapi.NewGoStarterClient(conn)
}
