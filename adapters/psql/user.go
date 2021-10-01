package psql

import (
	"context"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/domain/user"
	"gorm.io/gorm"
)

// Repository ...
type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create ...
func (r *Repository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Find ...
func (r *Repository) Find(ctx context.Context, conds ...interface{}) (*[]user.User, error) {
	var users []user.User

	if err := r.db.Find(&users, conds...).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// First ...
func (r *Repository) First(ctx context.Context, conds ...interface{}) (*user.User, error) {
	var user user.User

	if err := r.db.First(&user, conds...).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
