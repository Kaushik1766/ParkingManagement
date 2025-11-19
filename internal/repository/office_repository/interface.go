package officerepository

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/office_storage_mock.go -package=mocks
type OfficeStorage interface {
	AddOffice(ctx context.Context, officeName string, buildingID string, floorNumber int) error
	DeleteOffice(ctx context.Context, officeId string) error
	GetBuildingAndFloorByOffice(ctx context.Context, officeName string) (uuid.UUID, int, error)
	GetOfficesByBuilding(ctx context.Context, buildingID string) ([]models.Office, error)
	GetAllOffices(ctx context.Context) ([]models.Office, error)
	GetOfficeByName(ctx context.Context, officeName string) (models.Office, error)
}
