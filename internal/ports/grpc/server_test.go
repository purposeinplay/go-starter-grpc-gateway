package grpc_test

import (
	"context"
	"github.com/matryer/is"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/clientgrpc"
	portsgrpc "github.com/purposeinplay/go-starter-grpc-gateway/internal/ports/grpc"
	"go.uber.org/zap"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func newTestServer(
	t *testing.T,
	application app.Application,
) (
	*portsgrpc.Server,
	*grpc.ClientConn,
) {
	t.Helper()

	i := is.New(t)

	i.Helper()

	lis, dialer := newBufConnLis()

	logger, err := zap.NewDevelopment()
	i.NoErr(err)

	s := portsgrpc.NewGrpcTestServer(
		logger,
		application,
		lis,
	)
	i.NoErr(err)

	t.Cleanup(func() {
		_ = lis.Close()
		_ = s.Close()
	})

	c, err := clientgrpc.NewClientConn(
		"bufnet",
		clientgrpc.WithNoTLS(),
		clientgrpc.WithContextDialer(dialer),
	)
	i.NoErr(err)

	go func() {
		_ = s.ListenAndServe()
	}()

	return s, c
}

func newBufConnLis() (
	listener net.Listener,
	dialFunc func(context.Context, string) (net.Conn, error),
) {
	const bufSize = 1024 * 1024

	lis := bufconn.Listen(bufSize)

	return lis, func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}
