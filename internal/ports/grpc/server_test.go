package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql/psqltest"

	"github.com/purposeinplay/go-commons/logs"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"gorm.io/gorm"
)

const (
	dbUser     = "win"
	dbPassword = "pass"
	dbName     = "win"
)

var (
	port = "5432"
	dsn  = fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s h"+
			"ost=localhost "+
			"port=%s "+
			"sslmode=disable",
		dbUser,
		dbPassword,
		dbName,
		port,
	)
)

func newTestServer(
	t *testing.T,
) (*Server, *grpc.ClientConn, *gorm.DB) {
	t.Helper()

	ctx := context.Background()

	logger, err := logs.NewDevelopmentLogger()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	db, closeDB, err := psqltest.Connect(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	app := service.NewTestApplication(
		ctx,
		logger,
		cfg,
		db,
	)

	lis, dialer := newBufConnLis()

	srv := NewGrpcTestServer(ctx, logger, cfg, app, lis)

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			t.Error(err)
		}

		return
	}()

	conn, err := grpc.Dial(
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = closeDB()
		_ = conn.Close()
		_ = srv.Close()
	})

	return srv, conn, db
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
