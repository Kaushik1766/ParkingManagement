package authservice

import "github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"

type AuthenticationManager interface {
	Login(email, password string) (string, error)
	Signup(name, email, password string, role roles.Role) error
}
