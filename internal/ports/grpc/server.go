package grpc

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/purposeinplay/go-commons/auth"
	grpccommons "github.com/purposeinplay/go-commons/grpc"
	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server represents the GRPC server dependencies.
type Server struct {
	starterapi.UnimplementedGoStarterServer
	config     *config.Config
	logger     *zap.Logger
	jwtManager *auth.JWTManager
	server     *grpccommons.Server
	app        app.Application
}

// NewGrpcServer runs a grpc server.
func NewGrpcServer(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.Logger,
	application app.Application,
	jwtManager *auth.JWTManager,
) *Server {
	const servicePath = "/starter.apigrpc.GoStarter/"

	authRoles := map[string][]string{
		servicePath + "FindUser": {"user"},
	}

	authInterceptor := auth.NewAuthInterceptor(logger, jwtManager, authRoles)

	grpc_zap.ReplaceGrpcLoggerV2(logger)

	serverOptions := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			authInterceptor.Unary(),
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(
					grpc_ctxtags.CodeGenRequestFieldExtractor,
				),
			),
			grpc_zap.UnaryServerInterceptor(logger),
		),
	}

	srv := &Server{
		config:     cfg,
		logger:     logger,
		jwtManager: jwtManager,
		app:        application,
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
				fmt.Printf("cfg %+v", cfg)
				dialAddr := fmt.Sprintf(
					"127.0.0.1:%d",
					cfg.SERVER.Port-1,
				)
				if cfg.SERVER.Address != "" {
					dialAddr = fmt.Sprintf(
						"%v:%d",
						cfg.SERVER.Address,
						cfg.SERVER.Port-1,
					)
				}
				err := starterapi.RegisterGoStarterHandlerFromEndpoint(
					ctx,
					mux,
					dialAddr,
					dialOptions,
				)
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

// Stop terminates the server.
func (s Server) Stop() {
	s.server.Stop()
}

// Healthcheck endpoint.
func (Server) Healthcheck(
	context.Context,
	*emptypb.Empty,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
