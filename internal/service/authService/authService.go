package authservice

import userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/UserRepository"

type AuthService struct {
	db userrepository.UserStorage
}
