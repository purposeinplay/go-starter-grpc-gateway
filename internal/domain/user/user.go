package user

import (
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
)

// User domain model.
type User struct {
	id    uuid.UUID
	email string
}

// New instantiates a new user entity.
func New(
	id uuid.UUID,
	email string,
) (*User, error) {
	if id.IsZero() {
		return nil, errors.NewInvalidError("user id")
	}

	if email == "" {
		return nil, errors.NewEmailNotProvided()
	}

	return &User{
		id:    id,
		email: email,
	}, nil
}

// MustNew instantiates a new user entity.
// Panics if invalid data is given.
//
// This function is mainly used in testing.
func MustNew(
	id uuid.UUID,
	email string,
) *User {
	u, err := New(id, email)
	if err != nil {
		panic(err)
	}

	return u
}

// ID returns the user ID.
func (u User) ID() uuid.UUID {
	return u.id
}

// Email returns the user Email.
func (u User) Email() string {
	return u.email
}

// UnmarshalFromDatabase unmarshals User from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use it as a constructor - It may put domain into the invalid state!
func UnmarshalFromDatabase(
	id uuid.UUID,
	email string,
) *User {
	return &User{
		id:    id,
		email: email,
	}
}
