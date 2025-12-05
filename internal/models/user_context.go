package models

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"github.com/google/uuid"
)

type UserContext struct {
	Id    uuid.UUID
	Email string
	Role  roles.Role
}
