package ports

import (
	"context"
	"fmt"
	"github.com/purposeinplay/go-commons/auth"
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/service"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"testing"
	"time"

	"gorm.io/gorm"

	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
)


type StarterTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *grpcServer
	client starterapi.GoStarterClient
}

func (ts *StarterTestSuite) SetupTest() {
	ts.db.Exec("TRUNCATE TABLE users CASCADE")
}

func TestStarterTestSuite(t *testing.T) {
	cfg, err := config.LoadTestConfig("../../config.test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	logger, err := logs.NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	app, db := service.NewTestApplication(ctx, logger, cfg)

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.AccessTokenExp))

	srv := NewGrpcServer(ctx, cfg, logger, app, jwtManager)

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", cfg.SERVER.Address, cfg.SERVER.Port-1),
		grpc.WithInsecure(),
	)

	if err != nil {
		t.Fatal(err)
	}

	starterClient := starterapi.NewGoStarterClient(conn)

	ts := &StarterTestSuite{
		db:            db,
		server:        srv,
		client:        starterClient,
	}

	suite.Run(t, ts)
}
