package psql_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"

	// nolint: revive // reports line to long
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql/psqltest"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	var (
		ctx      = context.Background()
		i        = is.New(t)
		mockUser = psql.User{
			ID:    uuid.New(),
			Email: "test@test.com",
		}
	)

	insertMockUsers(t, dsn, mockUser)

	t.Run("AccountNotFound", func(t *testing.T) {
		i := i.New(t)

		db, closeDB, err := psqltest.NewDB(dsn)
		i.NoErr(err)

		t.Cleanup(func() {
			err := closeDB()
			i.NoErr(err)
		})

		r := psql.NewUserRepository(
			db,
		)

		us, err := r.GetUserByID(ctx, uuid.New())

		var appErr *errors.Error

		i.True(errors.As(err, &appErr))

		i.Equal(errors.ErrorTypeNotFound, appErr.Type())

		i.Equal(query.User{}, us)
	})

	t.Run("SuccessGetUserByID", func(t *testing.T) {
		i := i.New(t)

		db, closeDB, err := psqltest.NewDB(dsn)
		i.NoErr(err)

		t.Cleanup(func() {
			err := closeDB()
			i.NoErr(err)
		})

		r := psql.NewUserRepository(
			db,
		)

		us, err := r.GetUserByID(ctx, mockUser.ID)
		i.NoErr(err)

		i.Equal(query.User{
			ID:    mockUser.ID,
			Email: mockUser.Email,
		}, us)
	})

	t.Run("SuccessCreate", func(t *testing.T) {
		i := i.New(t)

		db, closeDB, err := psqltest.NewDB(dsn)
		i.NoErr(err)

		t.Cleanup(func() {
			err := closeDB()
			i.NoErr(err)
		})

		r := psql.NewUserRepository(
			db,
		)

		newUser := user.MustNew(
			uuid.New(),
			t.Name()+"@test.com",
		)

		err = r.CreateUser(ctx, newUser)
		i.NoErr(err)

		assertUserInDB(t, db, newUser)
	})
}

func insertMockUsers(t *testing.T, dsn string, users ...psql.User) {
	t.Helper()

	i := is.New(t)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		DSN:        dsn,
	}))
	i.NoErr(err)

	sqlDB, err := db.DB()
	i.NoErr(err)

	defer func() {
		err := sqlDB.Close()
		i.NoErr(err)
	}()

	err = db.Create(users).Error
	i.NoErr(err)
}

func assertUserInDB(t *testing.T, db *gorm.DB, acc *user.User) {
	t.Helper()

	i := is.New(t)

	var dbAcc *psql.User

	err := db.Take(&dbAcc, acc.ID()).Error
	i.NoErr(err)

	i.Equal(user.UnmarshalFromDatabase(
		dbAcc.ID,
		dbAcc.Email,
	),
		acc,
	)
}
