package query

import (
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// User represents the API model for the
// domain User.
type User struct {
	ID    uuid.UUID
	Email string
}
