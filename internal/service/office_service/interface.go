package officeservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
)

type OfficeMgr interface {
	AddOffice(ctx context.Context, officeName string, buildingName string, floorNumber int) error
	RemoveOffice(ctx context.Context, officeName string) error
	ListOfficesByBuilding(ctx context.Context, buildingName string) (map[int]string, error)
	GetAllOfficeNames(ctx context.Context) ([]string, error)
	GetOfficeByName(ctx context.Context, officeName string) (models.Office, error)
}
