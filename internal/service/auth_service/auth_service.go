package authservice

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
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

func (auth *AuthService) Signup(registerReq models.RegisterRequestDTO, role roles.Role) error {
	_, err := mail.ParseAddress(registerReq.Email)
	if err != nil {
		return errors.New("invalid email")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), 12)
	if err != nil {
		return err
	}

	_, err = auth.officeDb.GetOfficeByName(registerReq.Office)
	if role != roles.Admin && err != nil {
		return fmt.Errorf("error in signup service: %w", err)
	}

	err = auth.userDb.CreateUser(registerReq.Name, registerReq.Email, string(hashedPassword), registerReq.Office, role)
	return err
}

func (auth *AuthService) Login(loginReq models.LoginRequestDTO) (string, error) {
	_, err := mail.ParseAddress(loginReq.Email)
	if err != nil {
		return "", errors.New("invalid email")
	}
	user, err := auth.userDb.GetUserByEmail(loginReq.Email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		models.UserJwt{
			Email:  user.Email,
			Role:   user.Role,
			Office: user.Office.OfficeName,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        user.UserID.String(),
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
