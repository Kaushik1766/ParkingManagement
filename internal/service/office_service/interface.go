package officeservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/office_service_mock.go -package=mocks
type OfficeMgr interface {
	AddOffice(ctx context.Context, officeName string, buildingId string, floorNumber int) error
	RemoveOffice(ctx context.Context, officeId string) error
	ListOfficesByBuilding(ctx context.Context, buildingId string) ([]models.OfficeDTO, error)
	GetAllOfficeNames(ctx context.Context) ([]string, error)
	GetOfficeByName(ctx context.Context, officeName string) (models.Office, error)
}
