package auth

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

func ExtractTokenFromMetadata(md metadata.MD) (string, error) {
	meta, ok := md["authorization"]
	if !ok {
		return "", fmt.Errorf("auth token is missing")
	}

	auth := meta[0]

	const prefix = "Bearer "

	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("bad authorization string")
	}

	return auth[len(prefix):], nil
}
