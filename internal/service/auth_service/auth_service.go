package authservice

import (
	"context"
	"log"
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

func (auth *AuthService) Signup(ctx context.Context, registerReq models.RegisterRequestDTO, role roles.Role) error {
	_, err := mail.ParseAddress(registerReq.Email)
	if err != nil {
		log.Println(err)
		return customerrors.NewWebError(err, errorcodes.InvalidInput)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), 12)
	if err != nil {
		log.Println(err)
		return customerrors.NewWebError(err, errorcodes.InternalServerError)
	}

	err = auth.userDb.CreateUser(ctx, registerReq.Name, registerReq.Email, string(hashedPassword), registerReq.OfficeId, role)
	if err != nil {
		log.Println(err)
		return customerrors.NewWebError(err, errorcodes.UserAlreadyExists)
	}
	return nil
}

func (auth *AuthService) Login(ctx context.Context, loginReq models.LoginRequestDTO) (string, error) {
	_, err := mail.ParseAddress(loginReq.Email)
	if err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidInput)
	}
	user, err := auth.userDb.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidCredentials)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", customerrors.NewWebError(err, errorcodes.InvalidCredentials)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		models.UserJwt{
			Email:    user.Email,
			ID:       user.UserID.String(),
			Role:     user.Role,
			OfficeId: user.Office.OfficeID.String(),
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
