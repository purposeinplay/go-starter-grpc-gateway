package errors_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
)

func TestNewError(t *testing.T) {
	t.Parallel()

	internalCode := errors.InternalErrorCodeNotEnoughBalance

	t.Run("InvalidError", func(t *testing.T) {
		t.Parallel()

		t.Run("NoDetails", func(t *testing.T) {
			t.Parallel()

			err := errors.NewInvalidError("test")

			assertNoDetails(t, err)
		})

		t.Run("Details", func(t *testing.T) {
			t.Parallel()

			details := errors.NewDetails(internalCode, "test")

			err := errors.NewInvalidErrorWithDetails(
				"test",
				details,
			)

			assertDetails(t, err, details)
		})
	})

	t.Run("NotFoundError", func(t *testing.T) {
		t.Parallel()

		t.Run("NoDetails", func(t *testing.T) {
			t.Parallel()

			err := errors.NewNotFoundError("test")

			assertNoDetails(t, err)
		})

		t.Run("Details", func(t *testing.T) {
			t.Parallel()

			details := errors.NewDetails(internalCode, "test")

			err := errors.NewNotFoundErrorWithDetails(
				"test",
				details,
			)

			assertDetails(t, err, details)
		})
	})
}

func assertNoDetails(t *testing.T, err *errors.Error) {
	t.Helper()

	i := is.New(t)

	details, ok := err.Details()

	i.True(!ok)

	i.True(details == nil)
}

func assertDetails(
	t *testing.T,
	err *errors.Error,
	expectedDetails *errors.Details,
) {
	t.Helper()

	i := is.New(t)

	details, ok := err.Details()

	i.True(ok)

	i.True(details != nil)

	i.Equal(expectedDetails, details)
}

func TestIs(t *testing.T) {
	t.Run("EmptyDetailsSameErrorCodeDiffMessage", func(t *testing.T) {
		i := is.New(t)

		var (
			e1 = errors.NewInvalidError("err1")
			e2 = errors.NewInvalidError("err2")
		)

		i.True(errors.Is(e1, e2))
	})

	t.Run("EmptyDetailsDiffErrorCode", func(t *testing.T) {
		i := is.New(t)

		var (
			e1 = errors.NewInvalidError("err1")
			e2 = errors.NewNotFoundError("err2")
		)

		i.True(!errors.Is(e1, e2))
	})

	t.Run("SameErrorCodeSameInternalErrorCode", func(t *testing.T) {
		i := is.New(t)

		var (
			e1 = errors.NewInvalidErrorWithDetails(
				"err1",
				errors.NewDetails(
					errors.InternalErrorCodeNotEnoughBalance,
					"details msg",
				),
			)

			e2 = errors.NewInvalidErrorWithDetails(
				"err1",
				errors.NewDetails(
					errors.InternalErrorCodeNotEnoughBalance,
					"details msg",
				),
			)
		)

		i.True(errors.Is(e1, e2))
	})

	t.Run("SameErrorCodeDiffInternalErrorCode", func(t *testing.T) {
		i := is.New(t)

		var (
			e1 = errors.NewInvalidErrorWithDetails(
				"err1",
				errors.NewDetails(
					errors.InternalErrorCodeNotEnoughBalance,
					"details msg",
				),
			)

			e2 = errors.NewInvalidErrorWithDetails(
				"err1",
				errors.NewDetails(
					errors.InternalErrorCodeSessionExpired,
					"details msg",
				),
			)
		)

		i.True(!errors.Is(e1, e2))
	})
}
