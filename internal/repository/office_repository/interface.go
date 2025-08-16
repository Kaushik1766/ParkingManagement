package officerepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type OfficeStorage interface {
	AddOffice(officeName string, buildingID uuid.UUID, floorNumber int) error
	DeleteOffice(officeName string) error
	GetBuildingAndFloorByOffice(officeName string) (uuid.UUID, int, error)
	GetOfficesByBuilding(buildingID uuid.UUID) ([]models.Office, error)
	GetAllOffices() ([]models.Office, error)
	GetOfficeByName(officeName string) (models.Office, error)
}
