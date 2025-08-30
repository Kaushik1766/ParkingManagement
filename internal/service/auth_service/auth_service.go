package authservice

import (
	"net/mail"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	errorcodes "github.com/Kaushik1766/ParkingManagement/internal/constants/error_codes"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userDb userrepository.UserStorage
}

func NewAuthService(
	db userrepository.UserStorage,
) *AuthService {
	return &AuthService{
		userDb: db,
	}
}

func (auth *AuthService) Signup(registerReq models.RegisterRequestDTO, role roles.Role) error {
	_, err := mail.ParseAddress(registerReq.Email)
	if err != nil {
		return customerrors.NewWebError(err, errorcodes.InvalidInput)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), 12)
	if err != nil {
		return customerrors.NewWebError(err, errorcodes.InternalServerError)
	}

	err = auth.userDb.CreateUser(registerReq.Name, registerReq.Email, string(hashedPassword), registerReq.Office, role)
	if err != nil {
		return customerrors.NewWebError(err, errorcodes.UserAlreadyExists)
	}
	return nil
}

func (auth *AuthService) Login(loginReq models.LoginRequestDTO) (string, error) {
	_, err := mail.ParseAddress(loginReq.Email)
	if err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidInput)
	}
	user, err := auth.userDb.GetUserByEmail(loginReq.Email)
	if err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidCredentials)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidCredentials)
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
	signedToken, _ := jwtToken.SignedString([]byte(config.JWTSecret))
	return signedToken, nil
}
