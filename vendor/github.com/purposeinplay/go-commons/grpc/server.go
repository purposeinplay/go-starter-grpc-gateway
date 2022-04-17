package grpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ErrServerClosed indicates that the operation is now illegal because of
// the server has been closed.
var ErrServerClosed = errors.New("go-commons.grpc: server closed")

type (
	// the server interface defines basic methods for starting
	// and stopping a server.
	server interface {
		io.Closer
		Serve(listener net.Listener) error
	}

	// serverWithListener represents an implemented server
	// that can be started and stopped.
	serverWithListener struct {
		server       server
		listener     net.Listener
		running      bool
		runningMutex sync.Mutex
	}

	// registerServerFunc defines how we can register
	// a grpc service to a grpc server.
	registerServerFunc func(server *grpc.Server)

	// registerGatewayFunc defines how we can register
	// a grpc service to a gateway server.
	registerGatewayFunc func(
		mux *runtime.ServeMux,
		dialOptions []grpc.DialOption,
	) error
)

// ListenAndServe accepts incoming connections on the listener
// available in the struct.
func (s *serverWithListener) ListenAndServe() error {
	s.runningMutex.Lock()

	s.running = true

	s.runningMutex.Unlock()

	return s.server.Serve(s.listener)
}

// Close stops the seerver gracefully.
func (s *serverWithListener) Close() error {
	s.runningMutex.Lock()
	defer s.runningMutex.Unlock()

	if !s.running {
		return nil
	}

	err := s.server.Close()
	if err != nil {
		return err
	}

	return nil
}

type debugLogger interface {
	Debug(msg string, fields ...zap.Field)
}

// Server holds the grpc and gateway underlying servers.
// It starts and stops both of them together.
// In case one of the server fails the other one is closed.
type Server struct {
	grpcServerWithListener        *serverWithListener
	grpcGatewayServerWithListener *serverWithListener

	debug  bool
	logger debugLogger

	mu     sync.Mutex
	closed bool
}

// NewServer creates Server with both the grpc server and,
// if it's the case, the gateway server.
//
// The servers have not started to accept requests yet.
func NewServer(opt ...ServerOption) (*Server, error) {
	opts := defaultServerOptions()

	for _, o := range opt {
		o.apply(&opts)
	}

	aggregatorServer := new(Server)

	err := setDebugLogger(
		opts.debugLogger,
		aggregatorServer,
	)
	if err != nil {
		return nil, fmt.Errorf("set debugLogger: %w", err)
	}

	grpcServerWithListener, err := newGRPCServerWithListener(
		opts.grpcListener,
		opts.address,
		opts.tracing,
		opts.grpcServerOptions,
		opts.unaryServerInterceptors,
		opts.registerServer,
	)
	if err != nil {
		return nil, fmt.Errorf("new gRPC server: %w", err)
	}

	aggregatorServer.grpcServerWithListener = grpcServerWithListener

	// return here if a gateway server is not wanted
	if !opts.gateway {
		return aggregatorServer, nil
	}

	grpcGatewayServer, err := newGatewayServerWithListener(
		opts.muxOptions,
		opts.tracing,
		opts.registerGateway,
		opts.address,
		opts.httpMiddlewares,
	)
	if err != nil {
		return nil, fmt.Errorf("new gRPC gateway server: %w", err)
	}

	aggregatorServer.grpcGatewayServerWithListener = grpcGatewayServer

	return aggregatorServer, nil
}

func setDebugLogger(debugLogger debugLogger, server *Server) error {
	if debugLogger == nil {
		return nil
	}

	server.debug = true
	server.logger = debugLogger

	return nil
}

// ListenAndServe starts accepting incoming connections
// on both servers.
// If one of the servers encounters an error, both are stopped.
func (s *Server) ListenAndServe() error {
	s.mu.Lock()

	if s.closed {
		return ErrServerClosed
	}

	var g run.Group

	g.Add(
		s.runGRPCServer,
		func(err error) {
			_ = s.grpcServerWithListener.Close()
		},
	)

	// start gateway server.
	if s.grpcGatewayServerWithListener != nil {
		g.Add(
			s.runGatewayServer,
			func(err error) {
				_ = s.grpcGatewayServerWithListener.Close()
			},
		)
	}

	s.mu.Unlock()

	err := g.Run()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) runGRPCServer() error {
	s.logDebug(
		"starting gRPC server",
		zap.String(
			"address",
			s.grpcServerWithListener.listener.Addr().String(),
		),
	)

	return s.grpcServerWithListener.ListenAndServe()
}

func (s *Server) closeGRPCServer() error {
	s.logDebug("close grpc server")
	defer s.logDebug("grpc server done")

	return s.grpcServerWithListener.Close()
}

func (s *Server) runGatewayServer() error {
	s.logDebug(
		"starting gRPC gateway server for HTTP requests",
		zap.String(
			"address",
			s.grpcGatewayServerWithListener.listener.Addr().String(),
		),
	)

	return s.grpcGatewayServerWithListener.ListenAndServe()
}

func (s *Server) closeGatewayServer() error {
	s.logDebug("close gateway server")
	defer s.logDebug("gateway server done")

	return s.grpcGatewayServerWithListener.Close()
}

// Close closes both underlying servers.
// Safe to use concurrently and can be called multiple times.
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	// 1. Stop GRPC Gateway server first as it sits above GRPC server.
	// This also closes the underlying grpcListener.
	if s.grpcGatewayServerWithListener != nil {
		err := s.closeGatewayServer()
		if err != nil {
			return fmt.Errorf("close gateway server: %w", err)
		}
	}

	// 2. Stop GRPC server. This also closes the underlying grpcListener.
	err := s.closeGRPCServer()
	if err != nil {
		return fmt.Errorf("close grpc server: %w", err)
	}

	return nil
}

func (s *Server) logDebug(msg string, fields ...zap.Field) {
	if !s.debug {
		return
	}

	s.logger.Debug(msg, fields...)
}
