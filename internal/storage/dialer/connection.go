package dialer

import (
	// this is where we do the connections

	"fmt"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/config"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Connect will connect to that repository engine
func Connect(config *config.Config) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.DB.HOST, config.DB.USER, config.DB.PASSWORD, config.DB.NAME)
		conn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
		db = conn

		if err != nil {
			return errors.Wrap(err, "opening database connection")
		}

		return nil
	}

	err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5))

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

// Migrate runs the gorm migration for all models
func Migrate(db *gorm.DB, allModels []interface{}) error {
	if err := db.AutoMigrate(allModels...); err != nil {
		return err
	}

	return nil
}
