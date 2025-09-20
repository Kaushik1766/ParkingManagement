package models

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name     string     `gorm:"not null"`
	Email    string     `gorm:"not null;unique"`
	Password string     `gorm:"not null"`
	Role     roles.Role `gorm:"not null"`
	IsActive bool       `gorm:"default:true"`
	OfficeID uuid.UUID
	Office   Office
	Vehicles []Vehicle
}
type UserDTO struct {
	UserId string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Office string `json:"officeName"`
}

type UpdateUserDTO struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Office   string `json:"office,omitempty"`
}
type UserContext struct {
	Id    uuid.UUID
	Email string
	Role  roles.Role
}
