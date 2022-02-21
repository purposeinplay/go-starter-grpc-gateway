package psqlutil

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrProjectDirectoryNotFound is returned when the
// project directory is not found.
var ErrProjectDirectoryNotFound = errors.New("project directory not found")

func getProjectDirectoryPath(projectName string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", ErrUnableToResolveCaller
	}

	pathParts := strings.Split(filename, string(os.PathSeparator))

	// reverse range over path parts to find the walle directory
	// absolute path
	for len(pathParts) > 0 {
		p := pathParts[len(pathParts)-1]
		if p == projectName {
			return filepath.Join(pathParts...), nil
		}

		pathParts = pathParts[:len(pathParts)-1]
	}

	return "", ErrProjectDirectoryNotFound
}
