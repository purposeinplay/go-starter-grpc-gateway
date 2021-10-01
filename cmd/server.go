package cmd

import (
	"context"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/service"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/ports"

	"github.com/purposeinplay/go-commons/auth"
	"go.uber.org/zap"

	"github.com/purposeinplay/go-commons/logs"

	"github.com/spf13/cobra"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
)

var ServerCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

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

		jwtManager := auth.NewJWTManager(config.JWT.Secret, time.Duration(config.JWT.AccessTokenExp))

		app, cleanup := service.NewApplication(ctx, logger, config)
		defer func() {
			err := cleanup()
			if err != nil {
				logger.Fatal("error during cleanup %v", zap.Error(err))
			}
		}()

		server := ports.NewGrpcServer(ctx, config, logger, app, jwtManager)

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		logger.Info("Startup completed")
		<-c

		logger.Info("Shutdown started")
		server.Stop()
		logger.Info("Shutdown complete")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
