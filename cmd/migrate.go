package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/purposeinplay/go-commons/logs"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// migrateCmd subcommand that migrates the db.
var migrateCmd = &cobra.Command{
	Use: "migrate",
	Long: "Migrate database structures. " +
		"This will create new tables and add missing columns and indexes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := logs.NewLogger()
		if err != nil {
			return fmt.Errorf("new logger: %w", err)
		}

		defer func() {
			_ = logger.Sync()
		}()

		c, err := config.LoadConfig(cmd)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		db, err := psql.Connect(c)
		if err != nil {
			return fmt.Errorf("connect db: %w", err)
		}

		sql, err := db.DB()
		if err != nil {
			return fmt.Errorf("retrieve underlying db: %w", err)
		}

		driver, err := postgres.WithInstance(sql, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("new migrate postgres driver: %w", err)
		}

		migrator, err := migrate.NewWithDatabaseInstance(
			"file://./migrate/",
			c.DB.NAME, driver)
		if err != nil {
			return fmt.Errorf("new database migrate instance: %w", err)
		}

		err = migrator.Up()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("running migrations: %w", err)
			}

			logger.Info("could not find any new migrations", zap.Error(err))
		}

		version, dirty, err := migrator.Version()
		if err != nil {
			return fmt.Errorf("retrieve schema version: %w", err)
		}

		if dirty {
			return fmt.Errorf(
				"found dirty version: %w",
				&migrate.ErrDirty{Version: int(version)},
			)
		}

		logger.Info(
			"successfully run migrations command",
			zap.Uint("version", version),
		)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
