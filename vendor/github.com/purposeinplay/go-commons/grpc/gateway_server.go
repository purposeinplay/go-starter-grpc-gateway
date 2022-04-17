package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"google.golang.org/grpc"
)

var _ server = (*gatewayServer)(nil)

type gatewayServer struct {
	internalHTTPServer *http.Server
}

func (s *gatewayServer) Serve(listener net.Listener) error {
	err := s.internalHTTPServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *gatewayServer) Close() error {
	return s.internalHTTPServer.Shutdown(context.Background())
}

// nolint: revive // false-positive, it reports tracing as a control flag.
func newGatewayServerWithListener(
	muxOptions []runtime.ServeMuxOption,
	tracing bool,
	registerGateway registerGatewayFunc,
	address string,
	middlewares chi.Middlewares,
) (
	*serverWithListener,
	error,
) {
	grpcGatewayMux := runtime.NewServeMux(
		muxOptions...,
	)

	var handler http.Handler = grpcGatewayMux

	if tracing {
		handler = &ochttp.Handler{
			Handler:     grpcGatewayMux,
			Propagation: &propagation.HTTPFormat{},
		}
	}

	if registerGateway != nil {
		dialOptions := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		}

		err := registerGateway(grpcGatewayMux, dialOptions)
		if err != nil {
			return nil, fmt.Errorf("register gateway: %w", err)
		}
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("new listener: %w", err)
	}

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
	})

	r := chi.NewRouter()

	r.Use(middlewares...)
	r.Use(corsHandler.Handler)
	r.Get(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)
	r.Mount("/", handler)

	return &serverWithListener{
			server: &gatewayServer{
				internalHTTPServer: &http.Server{
					Handler: r,
				},
			},
			listener: listener,
		},
		nil
}
