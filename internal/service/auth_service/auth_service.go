package authservice

import (
	"fmt"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userDb   userrepository.UserStorage
	officeDb officerepository.OfficeStorage
}

func NewAuthService(
	db userrepository.UserStorage,
	officeDb officerepository.OfficeStorage,
) *AuthService {
	return &AuthService{
		userDb:   db,
		officeDb: officeDb,
	}
}

func (auth *AuthService) Signup(name, email, password, office string, role roles.Role) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = auth.officeDb.GetOfficeByName(office)
	if err != nil {
		return fmt.Errorf("error in signup service: %w", err)
	}

	err = auth.userDb.CreateUser(name, email, string(hashedPassword), office, role)
	return err
}

func (auth *AuthService) Login(email, password string) (string, error) {
	// _, err := mail.ParseAddress(email)
	// if err != nil {
	// 	return "", errors.New("invalid email")
	// }
	user, err := auth.userDb.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
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
	signedToken, err := jwtToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
