package officerepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/office_storage_mock.go -package=mocks
type OfficeStorage interface {
	AddOffice(officeName string, buildingID string, floorNumber int) error
	DeleteOffice(officeId string) error
	GetBuildingAndFloorByOffice(officeName string) (uuid.UUID, int, error)
	GetOfficesByBuilding(buildingID string) ([]models.Office, error)
	GetAllOffices() ([]models.Office, error)
	GetOfficeByName(officeName string) (models.Office, error)
}
