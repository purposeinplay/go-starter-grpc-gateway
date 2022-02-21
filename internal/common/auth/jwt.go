package auth

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-commons/auth"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"google.golang.org/grpc/metadata"
)

// UUIDFromContextJWT extracts the user id from the context.
func UUIDFromContextJWT(
	ctx context.Context,
	jwtManager *auth.JWTManager,
) (uuid.UUID, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	sign, err := auth.ExtractTokenFromMetadata(md)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("extract token from metadata: %w", err)
	}

	claims, err := jwtManager.Verify(sign)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("verify: %w", err)
	}

	parsedUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parse uuid: %w", err)
	}

	return parsedUUID, nil
}
