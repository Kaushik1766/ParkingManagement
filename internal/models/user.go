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
	OfficeID uuid.UUID  `gorm:"type:uuid;not null"`
	Office   Office     `gorm:"foreignKey:OfficeID;references:OfficeID"`
	Vehicles []Vehicle  `gorm:"foreignKey:UserID;references:UserID"`
}

func (u User) GetID() string {
	return u.UserID.String()
}
