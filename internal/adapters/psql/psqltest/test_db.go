// Package psqltest is similar package to net/http/http test
// used for storing testing utilities for testing the psql service
package psqltest

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	"github.com/goccy/go-yaml"
	"github.com/romanyx/polluter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// Driver used in order to open a new transaction when opnening a docker db.
	Driver = "pgsqltx"

	psqlDriver = "postgres"
)

// Pollute is a function to insert data in a datbase based on a YML.
// Used only for testing.
func Pollute(pollution interface{}, db *sql.DB) error {
	yamlPollution, err := yaml.Marshal(pollution)
	if err != nil {
		return fmt.Errorf("err while marshalling pollution yaml: %w", err)
	}

	return polluter.
		New(polluter.PostgresEngine(db)).
		Pollute(bytes.NewReader(yamlPollution))
}

// Register is a wrapper over txdb.Register.
func Register(dsn string) {
	txdb.Register(Driver, psqlDriver, dsn)
}

// NewDB returns a new transaction DB connection and inserts
// the given pollution data.
func NewDB(dsn string) (*gorm.DB, func() error, error) {
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DriverName: Driver,
				DSN:        dsn,
			}),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, nil, err
	}

	dbSQL, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, dbSQL.Close, nil
}
