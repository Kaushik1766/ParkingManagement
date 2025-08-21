package models

type RegisterRequestDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Office   string `json:"office" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
