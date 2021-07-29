package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/api"

	"github.com/purposeinplay/go-commons/auth"
	"go.uber.org/zap"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/storage/dialer"

	"github.com/purposeinplay/go-commons/logs"

	"github.com/spf13/cobra"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
)

var APICmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		config, err := config.LoadConfig(cmd)

		logger := logs.NewLogger()
		logger.Info("Go Starter API starting")

		if err != nil {
			logger.Fatal("unable to read config %v", zap.Error(err))
		}

		db, err := dialer.Connect(config)

		if err != nil {
			logger.Fatal("connecting to database: %+v", zap.Error(err))
		}

		jwtManager := auth.NewJWTManager(config.JWT.Secret, time.Duration(config.JWT.AccessTokenExp))

		api := api.NewAPI(ctx, config, logger, db, jwtManager)

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		logger.Info("Startup completed")
		<-c

		logger.Info("Shutdown started")
		api.Stop()
		logger.Info("Shutdown complete")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(APICmd)
}
