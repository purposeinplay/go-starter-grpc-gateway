package service

import (
	"context"
	psql2 "github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/command"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewApplication(ctx context.Context, logger *zap.Logger, config *config.Config) (app.Application, func() error) {
	db, err := psql2.Connect(config)

	if err != nil {
		logger.Fatal("connecting to database: %+v", zap.Error(err))
	}

	return newApplication(ctx, logger, config, db), func() error {
		return nil
	}
}

func newApplication(
	ctx context.Context,
	logger *zap.Logger,
	config *config.Config,
	db *gorm.DB,
) app.Application{

	userRepo := psql2.NewUserRepository(db)

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(logger, userRepo),
		},
		Queries:  app.Queries{
			FindUsers: query.NewFindUsersHandler(logger, userRepo),
			UserByID: query.NewUserById(logger, userRepo),
		},
	}
}

func NewTestApplication(
	ctx context.Context,
	logger *zap.Logger,
	config *config.Config,
) (app.Application, *gorm.DB) {
	db, err := psql2.Connect(config)

	if err != nil {
		logger.Fatal("connecting to database: %+v", zap.Error(err))
		panic(err)
	}

	return newApplication(ctx, logger, config, db), db
}