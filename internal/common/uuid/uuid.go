package uuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

// UUID represents a wrapper over the internal UUID implementation.
type UUID struct {
	internalUUID uuid.UUID
}

// New creates a new random uuid.
func New() UUID {
	return UUID{
		internalUUID: uuid.New(),
	}
}

// Parse decodes a string into an uuid.
func Parse(givenUUID string) (UUID, error) {
	parsedUUID, err := uuid.Parse(givenUUID)
	if err != nil {
		return UUID{}, err
	}

	return UUID{
		internalUUID: parsedUUID,
	}, nil
}

// MustParse decodes a string into an uuid.
// Panics if string is an invalid uuid.
func MustParse(givenUUID string) (UUID, error) {
	parsedUUID, err := Parse(givenUUID)
	if err != nil {
		panic(err)
	}

	return parsedUUID, nil
}

func (u UUID) String() string {
	return u.internalUUID.String()
}

// IsZero flags if the uuid is a zero value.
func (u UUID) IsZero() bool {
	return u.internalUUID == uuid.Nil
}

// Scan implements sql.Scanner and it's a wrapper
// over the internal uuid Scan Implementation.
func (u *UUID) Scan(src interface{}) error {
	return u.internalUUID.Scan(src)
}

// Value implements sql.Valuer and it's a wrapper
// over the internal uuid Value Implementation.
func (u UUID) Value() (driver.Value, error) {
	return u.internalUUID.Value()
}
