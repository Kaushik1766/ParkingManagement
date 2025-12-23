package models

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
	OfficeId string `json:"officeId,omitempty"`
}
