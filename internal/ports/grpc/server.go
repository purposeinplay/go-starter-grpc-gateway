package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"

	startergrpc "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc/v1"

	"github.com/purposeinplay/go-commons/auth"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/health/grpc_health_v1"

	grpccommons "github.com/purposeinplay/go-commons/grpc"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"go.uber.org/zap"
)

var _ startergrpc.GoStarterServer = (*Server)(nil)

// Server represents the GRPC server dependencies.
type Server struct {
	grpc_health_v1.UnimplementedHealthServer
	startergrpc.UnimplementedGoStarterServer

	ctx        context.Context
	logger     *zap.Logger
	app        app.Application
	cfg        *config.Config
	server     *grpccommons.Server
	jwtManager *auth.JWTManager
}

// NewGrpcServer runs a grpc server.
func NewGrpcServer(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	application app.Application,
	jwtManager *auth.JWTManager,
) *Server {
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	srv := &Server{
		ctx:        ctx,
		app:        application,
		cfg:        cfg,
		logger:     logger.Named("grpc.server"),
		jwtManager: jwtManager,
	}

	const servicePath = "/starter.apigrpc.GoStarter/"

	authRoles := map[string][]string{
		servicePath + "FindUser": {"user"},
	}

	authInterceptor := auth.NewAuthInterceptor(logger, jwtManager, authRoles)

	opts := []grpccommons.ServerOption{
		grpccommons.WithUnaryServerInterceptorRecovery(
			func(p any) (err error) {
				return srv.handlePanicRecover(p)
			},
		),

		grpccommons.WithUnaryServerInterceptor(
			authInterceptor.Unary(),
		),
		grpccommons.WithAddress(cfg.SERVER.Address),
		grpccommons.WithUnaryServerInterceptorLogger(
			logger.Named("grpc.server.interceptor"),
		),
		grpccommons.WithUnaryServerInterceptorCodeGen(),
		grpccommons.WithRegisterServerFunc(srv.registerGrpcServer),
		grpccommons.WithRegisterGatewayFunc(srv.registerGatewayServer),
		grpccommons.WithDebug(logger.Named("grpc.server.debug")),
	}

	grpcServer, err := grpccommons.NewServer(opts...)
	if err != nil {
		panic(err)
	}

	srv.server = grpcServer

	go func() {
		err := grpcServer.ListenAndServe()
		if err != nil {
			logger.Error("listen and serve err", zap.Error(err))
		}
	}()

	return srv
}

// NewGrpcTestServer returns a new grpc server to be used in tests.
func NewGrpcTestServer(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	application app.Application,
	listener net.Listener,
) *Server {
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	srv := &Server{
		ctx:    ctx,
		app:    application,
		cfg:    cfg,
		logger: logger.Named("grpc.server"),
	}

	opts := []grpccommons.ServerOption{
		// grpccommons.WithUnaryServerInterceptorRecovery(
		//	func(p interface{}) (err error) {
		//		return srv.handlePanicRecover(p)
		//	},
		// ),

		// grpccommons.WithUnaryServerInterceptorHandleErr(srv.handleErr),
		grpccommons.WithNoGateway(),
		grpccommons.WithDebug(logger),
		grpccommons.WithAddress(cfg.SERVER.Address),
		grpccommons.WithUnaryServerInterceptorLogger(
			logger.Named("grpc.test_server.interceptor"),
		),
		grpccommons.WithGRPCListener(listener),
		grpccommons.WithUnaryServerInterceptorCodeGen(),
		grpccommons.WithRegisterServerFunc(srv.registerGrpcServer),
		grpccommons.WithRegisterGatewayFunc(srv.registerGatewayServer),
		grpccommons.WithDebug(logger.Named("grpc.server.debug")),
	}

	grpcServer, err := grpccommons.NewServer(opts...)
	if err != nil {
		panic(err)
	}

	srv.server = grpcServer

	return srv
}

// ListenAndServe wraps the underlying
// grpccommons.Server ListenAndServe method.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Close terminates the server.
func (s *Server) Close() error {
	return s.server.Close()
}

// Healthcheck endpoint.
func (*Server) Healthcheck(
	context.Context,
	*emptypb.Empty,
) (*emptypb.Empty, error) {
	return nil, nil
}

// Check is used to verify if the API is able to accept requests.
func (*Server) Check(
	context.Context,
	*grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch works like Check, but it creates an ongoing connection.
func (*Server) Watch(
	*grpc_health_v1.HealthCheckRequest,
	grpc_health_v1.Health_WatchServer,
) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *Server) registerGrpcServer(server *grpc.Server) {
	startergrpc.RegisterGoStarterServer(server, s)
	grpc_health_v1.RegisterHealthServer(server, s)
}

func (s *Server) registerGatewayServer(
	mux *runtime.ServeMux,
	dialOptions []grpc.DialOption,
) error {
	host, port, err := parseHostPort(s.cfg.SERVER.Address)
	if err != nil {
		return fmt.Errorf("could not parse server address: %w", err)
	}

	err = startergrpc.RegisterGoStarterHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%v:%v", host, port-1),
		dialOptions,
	)
	if err != nil {
		return fmt.Errorf("register gRPC gateway: %w", err)
	}

	return nil
}

func parseHostPort(address string) (string, int, error) { //nolint:gocritic
	hostString, portString, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, fmt.Errorf("invalid address: %w", err)
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return "", 0, fmt.Errorf("parse port: %w", err)
	}

	return hostString, port, nil
}
