package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/purposeinplay/go-commons/auth"
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/ports/grpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/service"
	"github.com/spf13/cobra"
)

// ServerCmd subcommand that starts the server.
var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, _ := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			os.Kill,
		)

		logger, err := logs.NewLogger()
		defer func() {
			_ = logger.Sync()
		}()

		if err != nil {
			return fmt.Errorf("could not create logger %w", err)
		}

		config, err := config.LoadConfig(cmd)
		if err != nil {
			return fmt.Errorf("unable to read config %w", err)
		}

		logger.Info("Openmatch API starting")

		jwtManager := auth.NewJWTManager(
			config.JWT.Secret,
			time.Duration(config.JWT.AccessTokenExp),
		)

		app, cleanup := service.NewApplication(ctx, logger, config)
		defer func() {
			err := cleanup()
			if err != nil {
				logger.Fatal("error during cleanup", zap.Error(err))
			}
		}()

		server := grpc.NewGrpcServer(ctx, logger, config, app, jwtManager)
		defer func() {
			err := server.Close()
			if err != nil {
				logger.Fatal("close grpc server", zap.Error(err))
			}
		}()
		logger.Info("Startup completed")

		<-ctx.Done()

		logger.Info("Shutdown complete")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
