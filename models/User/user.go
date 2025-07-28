package user

import (
	"github.com/Kaushik1766/ParkingManagement/models/enums"
	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID
	Name     string
	Email    string
	Password string
	Role     enums.Role
}

func (u User) GetId() string {
	return u.UserId.String()
}
