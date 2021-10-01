package ports

import (
	"context"
	"fmt"
	grpccommons "github.com/purposeinplay/go-commons/grpc"
	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/purposeinplay/go-commons/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
)

type grpcServer struct {
	starterapi.UnimplementedGoStarterServer
	config     *config.Config
	logger     *zap.Logger
	jwtManager *auth.JWTManager
	server *grpccommons.Server
	app app.Application
}

func NewGrpcServer(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	app app.Application,
	jwtManager *auth.JWTManager,
) *grpcServer {

	const servicePath = "/starter.apigrpc.GoStarter/"
	authRoles := map[string][]string{
		servicePath + "FindUser": {"user"},
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

	srv := &grpcServer{
		app: app,
		config:     config,
		logger:     logger,
		jwtManager: jwtManager,
	}

	options := []grpccommons.ServerOption{
		grpccommons.WithGrpcServerOptions(serverOptions),

		grpccommons.ReplaceLogger(logger),

		grpccommons.RegisterServer(
			func(server *grpc.Server) {
				starterapi.RegisterGoStarterServer(server, srv)
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
	srv.server = grpcServer

	return srv
}

func (a *grpcServer) Stop() {
	a.server.Stop()
}

func (a *grpcServer) Healthcheck(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
