package user

import "context"

type UserRepository interface {
	Find(ctx context.Context, conds ...interface{}) (*[]User, error)
	First(ctx context.Context, conds ...interface{}) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}