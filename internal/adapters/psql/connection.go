package psql

import (

	// this is where we do the connections.
	"fmt"

	"gorm.io/driver/postgres"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect will connect to that repository engine.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		url := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.HOST,
			cfg.DB.USER,
			cfg.DB.PASSWORD,
			cfg.DB.NAME,
		)
		conn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
		db = conn

		if err != nil {
			return errors.Wrap(err, "opening database connection")
		}

		return nil
	}

	const maxRetries = 5

	err := backoff.Retry(
		operation,
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries),
	)
	if err != nil {
		return nil, err
	}

	db.Logger = db.Logger.LogMode(logger.Info)

	sqlDB, err := db.DB()

	if err == nil {
		err = sqlDB.Ping()
	}

	if err != nil {
		return nil, errors.Wrap(err, "checking database connection")
	}

	return db, nil
}

// Migrate runs the gorm migration for all models.
func Migrate(db *gorm.DB, allModels []interface{}) error {
	if err := db.AutoMigrate(allModels...); err != nil {
		return err
	}

	return nil
}
