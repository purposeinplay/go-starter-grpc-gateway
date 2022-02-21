package psqlutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrUnableToResolveCaller is returned when the caller CWD cannot be retrieved.
var ErrUnableToResolveCaller = errors.New("unable to resolve caller")

// ReadSchema reads schema dynamically based on the CWD of the caller.
func ReadSchema() (string, error) {
	p, err := getProjectDirectoryPath("go-starter-grpc-gateway")
	if err != nil {
		return "", fmt.Errorf("get project directory: %w", err)
	}

	schemaPath := filepath.Clean(
		filepath.Join(
			string(os.PathSeparator),
			filepath.Join(
				p,
				"sql",
				"schema.sql",
			),
		),
	)

	schemaB, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("err while reading schema: %w", err)
	}

	return string(schemaB), nil
}
