package userservice

import (
	"context"
)

// will pass primaryKey in the context

type UserManager interface {
	UpdateProfile(ctx context.Context, name, email, password string) error
	DeleteProfile(ctx context.Context) error
}
