package service

import (
	"context"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewApplication returns a production application.
func NewApplication(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
) (
	application app.Application,
	cleanup func() error,
) {
	db, err := psql.Connect(cfg)
	if err != nil {
		logger.Fatal("connecting to database: %+v", zap.Error(err))
	}

	return newApplication(
			ctx,
			logger,
			cfg,
			db,
		), func() error {
			return nil
		}
}

func newApplication(
	_ context.Context,
	_ *zap.Logger,
	_ *config.Config,
	db *gorm.DB,
) app.Application {
	userRepo := psql.NewUserRepository(db)

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.MustNewCreateUserHandler(userRepo),
		},
		Queries: app.Queries{
			FindUsers: query.MustNewFindUsersHandler(userRepo),
			UserByID:  query.MustNewUserByIDHandler(userRepo),
		},
	}
}

// NewTestApplication instantiates a test application.
func NewTestApplication(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	db *gorm.DB,
) app.Application {
	return newApplication(ctx, logger, cfg, db)
}
