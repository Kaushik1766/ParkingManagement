package user

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/google/uuid"
)

type UserContext struct {
	Id    uuid.UUID
	Email string
	Role  enums.Role
}
