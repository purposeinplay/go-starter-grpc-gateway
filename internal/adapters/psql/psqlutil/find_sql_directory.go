package psqlutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindSQLDir attempts to compose the path to the sql directory in the project.
func FindSQLDir() (string, error) {
	p, err := getProjectDirectoryPath("go-starter-grpc-gateway")
	if err != nil {
		return "", fmt.Errorf("get project directory: %w", err)
	}

	sqlPath := filepath.Clean(
		filepath.Join(
			string(os.PathSeparator),
			filepath.Join(
				p,
				"sql",
			),
		),
	)

	_, err = os.Stat(
		sqlPath,
	)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		}

		return "", fmt.Errorf("err while checking if sql dir exists: %w", err)
	}

	return sqlPath, nil
}
