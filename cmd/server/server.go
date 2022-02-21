package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/purposeinplay/go-commons/auth"
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/ports/grpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// CMD represents the http server command.
var CMD = &cobra.Command{
	Use: "http",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, _ := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			os.Kill,
		)

		logger, err := logs.NewLogger()
		if err != nil {
			return fmt.Errorf("new logger: %w", err)
		}

		defer func() {
			_ = logger.Sync()
		}()

		config, err := config.LoadConfig(cmd)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		jwtManager := auth.NewJWTManager(
			config.JWT.Secret,
			time.Duration(config.JWT.AccessTokenExp),
		)

		app, cleanupApp := service.NewApplication(ctx, logger, config)

		defer func() {
			err := cleanupApp()
			if err != nil {
				logger.Error("application cleanup", zap.Error(err))
			}
		}()

		server := grpc.NewGrpcServer(ctx, config, logger, app, jwtManager)
		defer server.Stop()

		logger.Info("Startup completed")

		<-ctx.Done()

		logger.Info("Shutdown complete")

		return nil
	},
}
