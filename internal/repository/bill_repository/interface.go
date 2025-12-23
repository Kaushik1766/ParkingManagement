package billrepository

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/bill_repository_mock.go -package=mocks
type BillStorage interface {
	SaveBill(ctx context.Context, bill models.BillDTO) error
	GetBill(ctx context.Context, userId string, month, year int) (models.BillDTO, error)
}
