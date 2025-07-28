package user

import "github.com/google/uuid"

type Role int

const (
	Customer Role = iota
	Admin
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "Admin"
	default:
		return "Customer"
	}
}

type User struct {
	UserId   uuid.UUID
	Name     string
	Email    string
	Password string
	Role     Role
}

func (u User) GetId() string {
	return u.UserId.String()
}
