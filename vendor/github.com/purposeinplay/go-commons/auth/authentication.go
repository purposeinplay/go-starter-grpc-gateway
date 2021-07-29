package auth

import (
	"context"
	"log"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type AuthInterceptor struct {
	logger     *zap.Logger
	jwtManager *JWTManager
	authRoles  map[string][]string
}

func NewAuthInterceptor(logger *zap.Logger, jwtManager *JWTManager, authRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		logger:     logger,
		jwtManager: jwtManager,
		authRoles:  authRoles,
	}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		ctxAuth, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctxAuth, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	authRoles, ok := i.authRoles[method]
	if !ok {
		// public route
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		i.logger.Error("metadata was not provided")
		return ctx, status.Errorf(codes.Unauthenticated, "auth token is invalid")
	}

	signature, err := ExtractTokenFromMetadata(md)

	if err != nil {
		i.logger.Error("ExtractTokenFromMetadata error", zap.Error(err))
		return ctx, status.Errorf(codes.Unauthenticated, "auth token is invalid")
	}

	claims, err := i.jwtManager.Verify(signature)

	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "auth token is invalid")
	}

	authCtx := WithUser(ctx, claims)

	for _, role := range authRoles {
		if role == claims.Role {
			return authCtx, nil
		}
	}

	return ctx, status.Error(codes.Unauthenticated, "no permission to access this RPC")
}
