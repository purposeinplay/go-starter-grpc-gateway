package cmd

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing columns and indexes.",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logs.NewLogger()

		if err != nil {
			log.Panicf("could not create logger %+v", err)
		}

		defer func() {
			_ = logger.Sync()
		}()

		c, err := config.LoadConfig(cmd)

		if err != nil {
			logger.Fatal("Unable to read config", zap.Error(err))
		}

		db, err := psql.Connect(c)

		if err != nil {
			logger.Fatal("error opening database", zap.Error(err))
		}

		sql, err := db.DB()

		if err != nil {
			logger.Fatal("invalid db", zap.Error(err))
		}

		driver, err := postgres.WithInstance(sql, &postgres.Config{})

		if err != nil {
			logger.Fatal("could not get instance", zap.Error(err))
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://./migrate/",
			c.DB.NAME, driver)

		if err != nil {
			logger.Fatal("could not retrieve db migrator instance", zap.Error(err))
		}

		if err := migrator.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Info("could not find any new migrations", zap.Error(err))
				return
			}
			logger.Fatal("error while performing migrations", zap.Error(err))
		}

		version, _, err := migrator.Version()
		if err != nil {
			logger.Fatal("could not get migration version", zap.Error(err))
		}

		logger.Info("successfully migrated DB:", zap.Uint("version", version))
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}