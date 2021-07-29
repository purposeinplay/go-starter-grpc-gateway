package user

import (
	"errors"
	starterapi "github.com/purposeinplay/go-starter-grpc-gateway/apigrpc"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain"
)

// User model
type User struct {
	domain.Base
	Email string 				`json:"email" gorm:"uniqueIndex; not null"`
}

type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

func New(text string) error {
	return &Error{text}
}

var (
	ErrUserNotFound   = New("user not found")
	ErrEmailIsRequired   = New("email is required")
	ErrCouldNotCreateUser   = New("could not create user")
)

func GRPCErrorFromError(err error) starterapi.ErrorCode_Code {
	var starterErr starterapi.ErrorCode_Code

	if errors.Is(err, ErrEmailIsRequired) {
		starterErr = starterapi.ErrorCode_EMAIL_REQUIRED_ERROR
	} else {
		starterErr = starterapi.ErrorCode_TYPE_UNSPECIFIED
	}

	return starterErr
}