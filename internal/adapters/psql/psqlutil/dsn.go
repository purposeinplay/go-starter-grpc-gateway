package psqlutil

import (
	"fmt"
)

// ComposeDSN returns a PostgreSQL Data Source Name.
func ComposeDSN(
	host,
	port,
	user,
	password,
	dbName,
	sslMode string,
) string {
	return fmt.Sprintf(
		"host=%s "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"port=%s "+
			"sslmode=%s",
		host,
		user,
		password,
		dbName,
		port,
		sslMode,
	)
}
