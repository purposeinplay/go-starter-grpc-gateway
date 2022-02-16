package errors

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrorType holds the canonical status codes for GRPC and HTTP.
//
// The error types defined are meant to be generic.
// They are mapped to the status codes available in HTTP (404,401)
// and grpc codes package("google.golang.org/grpc/codes").
// NotFound, BadRequest...
type ErrorType struct {
	t string
}

func (e ErrorType) String() string {
	return e.t
}

var (
	// ErrorTypeInvalid is used when invalid data is
	// passed to the application.
	// Maps to:
	// HTTP: 400
	// GRPC: 3.
	ErrorTypeInvalid = ErrorType{"invalid"}

	// ErrorTypeNotFound is used when a resource is not found.
	// Maps to:
	// HTTP: 404
	// GRPC: 5.
	ErrorTypeNotFound = ErrorType{"not-found"}

	// ErrorTypeInternal is used when an internal error
	// is thrown by the system
	// Maps to:
	// HTTP: 500
	// GRPC: 13.
	ErrorTypeInternal = ErrorType{"internal"}

	// ErrorTypeUnauthorized is used when an user attempts
	// to perform an unauthorized action
	// Maps to:
	// HTTP: 401
	// GRPC: 7.
	ErrorTypeUnauthorized = ErrorType{"unauthorized"}
)

// ApplicationErrorCode holds error codes specific to the application.
// They are also defined in the proto files.
type ApplicationErrorCode struct {
	c uint8
}

func (e ApplicationErrorCode) String() string {
	return strconv.Itoa(int(e.c))
}

var (
	// InternalErrorCodeNotEnoughBalance is used when the user balance
	// is not enough to perform a transaction.
	InternalErrorCodeNotEnoughBalance = ApplicationErrorCode{c: 1}

	// InternalErrorCodeSessionExpired is returned when a user attempts
	// to perform an authenticated action but their session is expired.
	InternalErrorCodeSessionExpired = ApplicationErrorCode{c: 2}

	// InternalErrorCodeEmailNotProvided is when a user attempts to register
	// but does not provide an email.
	InternalErrorCodeEmailNotProvided = ApplicationErrorCode{c: 3}
)

// Error represents an application error.
type Error struct {
	// Canonical status code.
	t ErrorType

	// Human-readable error message.
	msg string

	// Details holds extra information about the error.
	details *Details
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"application error: type: %s message: %s",
		e.t.String(),
		e.msg,
	)
}

// Type returns the error type (canonical status).
func (e *Error) Type() ErrorType {
	return e.t
}

// Is implements the errors.Is interface.
func (e *Error) Is(target error) bool {
	var (
		err *Error
		sameErrorCode,
		sameInternalErrorCode bool
	)

	if errors.As(target, &err) {
		sameErrorCode = err.t == e.t

		if sameErrorCode {
			var internalErr1,
				internalErr2 ApplicationErrorCode

			if err.details != nil {
				internalErr1 = err.details.applicationErrorCode
			}

			if e.details != nil {
				internalErr2 = e.details.applicationErrorCode
			}

			sameInternalErrorCode = internalErr1 == internalErr2
		}

		return sameErrorCode && sameInternalErrorCode
	}

	return false
}

// Details returns the error details and a flag that
// set to true if those details are available.
func (e *Error) Details() (*Details, bool) {
	if e.details == nil {
		return nil, false
	}

	return e.details, true
}

// Message returns a high level description of the error.
func (e *Error) Message() string {
	return e.msg
}

// Details represents extra data attached to an application error.
type Details struct {
	applicationErrorCode ApplicationErrorCode
	message              string
}

// ApplicationErrorCode returns the application error code contained
// in the Details struct.
func (d *Details) ApplicationErrorCode() ApplicationErrorCode {
	return d.applicationErrorCode
}

// Message returns the data from the message field.
func (d *Details) Message() string {
	return d.message
}

// NewDetails creates a new Details struct.
func NewDetails(
	internalCode ApplicationErrorCode,
	msg string,
) *Details {
	return &Details{
		applicationErrorCode: internalCode,
		message:              msg,
	}
}

// NewInvalidError creates a new application Invalid Error.
func NewInvalidError(msg string) *Error {
	return &Error{
		t:   ErrorTypeInvalid,
		msg: msg,
	}
}

// NewUnauthorizedError creates a new application unauthorized error.
func NewUnauthorizedError(msg string) *Error {
	return &Error{
		t:   ErrorTypeUnauthorized,
		msg: msg,
	}
}

// NewInvalidErrorWithDetails creates a new application Invalid Error
// and attaches a Details object to it.
func NewInvalidErrorWithDetails(
	msg string,
	details *Details,
) *Error {
	err := NewInvalidError(msg)

	err.details = details

	return err
}

// NewNotFoundError creates a new application Not Found Error.
func NewNotFoundError(msg string) *Error {
	return &Error{
		t:   ErrorTypeNotFound,
		msg: msg,
	}
}

// NewNotFoundErrorWithDetails creates a new application Not Found Error.
// and attaches a Details object to it.
func NewNotFoundErrorWithDetails(
	msg string,
	details *Details,
) *Error {
	err := NewNotFoundError(msg)

	err.details = details

	return err
}

// NewInternalError creates a new application Internal Error.
// This type of errors should never be shown to a user.
func NewInternalError(msg string) *Error {
	return &Error{
		t:       ErrorTypeInvalid,
		msg:     msg,
		details: nil,
	}
}

// NewEmailNotProvided create a new Application Invalid Error
// with an application error code specifying an email not proivede err.
func NewEmailNotProvided() *Error {
	return &Error{
		t:   ErrorTypeInvalid,
		msg: "email not provided",
		details: &Details{
			applicationErrorCode: InternalErrorCodeEmailNotProvided,
			message:              "",
		},
	}
}

// As wraps errors.As.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is wraps errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
