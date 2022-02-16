package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/purposeinplay/go-commons/auth"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"

	// nolint: revive // reports line to long
	portsgrpc "github.com/purposeinplay/go-starter-grpc-gateway/internal/ports/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//
// type StarterTestSuite struct {
// 	suite.Suite
// 	db     *gorm.DB
// 	Server *ports.Server
// 	client starterapi.GoStarterClient
// }
//
// func (ts *StarterTestSuite) SetupTest() {
// 	ts.db.Exec("TRUNCATE TABLE users CASCADE")
// }
//
// func TestStarterTestSuite(t *testing.T) {
// 	cfg, err := config.LoadTestConfig("../../config.test.yaml")
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	ctx := context.Background()
// 	logger, err := logs.NewLogger()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer func() {
// 		_ = logger.Sync()
// 	}()
//
// 	app, db := service.NewTestApplication(ctx, logger, cfg)
//
// 	jwtManager := auth.NewJWTManager(cfg.JWT.Secret,
// 		time.Duration(cfg.JWT.AccessTokenExp))
//
// 	srv := grpc.NewGrpcServer(ctx, cfg, logger, app, jwtManager)
//
// 	conn, err := grpc.Dial(
// 		fmt.Sprintf("%v:%v", cfg.SERVER.Address, cfg.SERVER.Port-1),
// 		grpc.WithInsecure(),
// 	)
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	starterClient := starterapi.NewGoStarterClient(conn)
//
// 	ts := &StarterTestSuite{
// 		db:     db,
// 		Server: srv,
// 		client: starterClient,
// 	}
//
// 	suite.Run(t, ts)
// }

func newServerOnRandomPort(
	t *testing.T,
	application app.Application,
	cfg *config.Config,
) (
	*portsgrpc.Server,
	*grpc.ClientConn,
) {
	t.Helper()

	i := is.New(t)

	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.AccessTokenExp),
	)

	s := portsgrpc.NewGrpcServer(
		context.TODO(),
		cfg,
		zap.NewExample(),
		application,
		jwtManager,
	)

	c, err := grpc.Dial(
		fmt.Sprintf(
			"127.0.0.1:%d",
			cfg.SERVER.Port-1,
		),
		grpc.WithInsecure(),
	)
	i.NoErr(err)

	return s, c
}
