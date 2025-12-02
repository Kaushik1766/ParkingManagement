package userrepository

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/user_storage_mock.go -package=mocks
type UserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserById(ctx context.Context, id string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	Save(ctx context.Context, user models.User) error
	CreateUser(ctx context.Context, name, email, password, officeId string, role roles.Role) error
}
