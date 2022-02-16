package psql

import (
	"context"
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/app/query"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/errors"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/uuid"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"gorm.io/gorm"
)

var (
	_ user.Repository          = (*Repository)(nil)
	_ query.UserByIDReadModel  = (*Repository)(nil)
	_ query.FindUsersReadModel = (*Repository)(nil)
)

// User represents the user model in the PostgreSQL database.
type User struct {
	ID    uuid.UUID `validate:"required" gorm:"primaryKey;column:user_id"`
	Email string    `validate:"required,email"`
}

// TableName satisfies the gorm.Tabler interface.
func (User) TableName() string {
	return "users"
}

// Repository represents a PostgreSQL User Repository.
type Repository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL User Repository.
func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the PostgreSQL database.
func (r Repository) CreateUser(ctx context.Context, u *user.User) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		psqlUser, err := r.marshalUser(u)
		if err != nil {
			return fmt.Errorf("marshal user: %w", err)
		}

		err = createUser(ctx, tx, psqlUser)
		if err != nil {
			return fmt.Errorf("create user query: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("tx sql: %w", err)
	}

	return nil
}

// FindUsers queries the PostgreSQL database for users.
func (r Repository) FindUsers(
	ctx context.Context,
	filter user.Filter,
) ([]query.User, error) {
	var users []query.User

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sqlUsers, err := findUsers(ctx, tx, filter)
		if err != nil {
			return fmt.Errorf("find users query: %w", err)
		}

		users = make([]query.User, 0, len(sqlUsers))

		for _, u := range sqlUsers {
			users = append(users, query.User{
				ID:    u.ID,
				Email: u.Email,
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("tx sql: %w", err)
	}

	return users, nil
}

func (Repository) marshalUser(u *user.User) (*User, error) {
	psqlUser := &User{
		ID:    u.ID(),
		Email: u.Email(),
	}

	err := validate.Struct(psqlUser)
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return psqlUser, nil
}

// GetUserByID queries the PostgreSQL database for a
// user with the given ID.
func (r Repository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (query.User, error) {
	var u query.User

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sqlUsers, err := findUsers(ctx, tx, user.Filter{ID: id})
		if err != nil {
			return fmt.Errorf("find users query: %w", err)
		}

		if len(sqlUsers) == 0 {
			return errors.NewNotFoundError("user")
		}

		u = query.User{
			ID:    sqlUsers[0].ID,
			Email: sqlUsers[0].Email,
		}

		return nil
	})
	if err != nil {
		return query.User{}, fmt.Errorf("tx sql: %w", err)
	}

	return u, nil
}

func createUser(
	ctx context.Context,
	db *gorm.DB,
	u *User,
) error {
	err := db.WithContext(ctx).Create(u).Error
	if err != nil {
		return fmt.Errorf("execute create user query: %w", err)
	}

	return nil
}

func findUsers(
	ctx context.Context,
	db *gorm.DB,
	filter user.Filter,
) ([]*User, error) {
	session := db.WithContext(ctx)

	if !filter.ID.IsZero() {
		session = session.Where("user_id = ?", filter.ID.String())
	}

	if filter.Email != nil {
		session = session.Where("email = ?", *filter.Email)
	}

	var users []*User

	err := session.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("execut find users query: %w", err)
	}

	return users, nil
}
