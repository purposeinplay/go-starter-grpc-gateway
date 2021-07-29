package cmd

import (
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/storage/dialer"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing columns and indexes.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logs.NewLogger()

		c, err := config.LoadConfig(cmd)

		if err != nil {
			logger.Fatal("Unable to read config", zap.Error(err))
		}

		db, err := dialer.Connect(c)

		if err != nil {
			logger.Fatal("error opening database", zap.Error(err))
		}

		models := []interface{}{
			&user.User{},
		}

		if err = dialer.Migrate(db, models); err != nil {
			logger.Fatal("error while performing migration", zap.Error(err))
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
