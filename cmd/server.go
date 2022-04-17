package cmd

import (
	"context"
	"log"
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

var ServerCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
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
			log.Panicf("could not create logger %+v", err)
		}

		config, err := config.LoadConfig(cmd)
		if err != nil {
			logger.Fatal("unable to read config %v", zap.Error(err))
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
				logger.Fatal("error during cleanup %v", zap.Error(err))
			}
		}()

		server := grpc.NewGrpcServer(ctx, logger, config, app, jwtManager)
		defer func() {
			err := server.Close()
			if err != nil {
				logger.Error("close grpc server", zap.Error(err))
			}
		}()
		logger.Info("Startup completed")

		<-ctx.Done()

		logger.Info("Shutdown complete")

		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
