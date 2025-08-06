package officerepository

import "github.com/Kaushik1766/ParkingManagement/internal/models/office"

type OfficeStorage interface {
	AddOffice(officeName, buildingName string, floorNumber int) error
	DeleteOffice(officeName string) error
	GetBuildingAndFloorByOffice(officeName string) (string, int, error)
	GetOfficesByBuilding(buildingName string) ([]office.Office, error)
	GetAllOffices() ([]office.Office, error)
	GetOfficeByName(officeName string) (office.Office, error)
}
