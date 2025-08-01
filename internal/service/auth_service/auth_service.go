package authservice

import (
	"errors"
	"net/mail"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db userrepository.UserStorage
}

func NewAuthService(db userrepository.UserStorage) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (auth *AuthService) Signup(name, email, password string, role enums.Role) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = auth.db.CreateUser(name, email, string(hashedPassword), role)
	return err
}

func (auth *AuthService) Login(email, password string) (*jwt.Token, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return &jwt.Token{}, errors.New("invalid email")
	}
	user, err := auth.db.GetUserByEmail(email)
	if err != nil {
		return &jwt.Token{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return &jwt.Token{}, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		userjwt.UserJwt{
			Email: user.Email,
			ID:    user.UserId.String(),
			Role:  user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	)
	return jwtToken, nil
}
