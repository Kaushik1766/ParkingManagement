package userservice

import (
	"context"

	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo userrepository.UserStorage
}

func NewUserService(repo userrepository.UserStorage) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) UpdateProfile(ctx context.Context, name, email, password string) error {
	currentUser := ctx.Value("user").(user.User)
	currentUser.Name = name
	currentUser.Email = email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	currentUser.Password = string(hashedPassword)
	err = s.userRepo.Save(currentUser)
	return err
}

func (s *UserService) DeleteProfile(ctx context.Context) error {
	currentUser := ctx.Value("user").(user.User)
	currentUser.IsActive = false
	err := s.userRepo.Save(currentUser)
	return err
}
