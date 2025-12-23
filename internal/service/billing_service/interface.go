package billingservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/billing_service_mock.go -package=mocks
type BillingMgr interface {
	GetMonthlyBill(ctx context.Context, userId string, month, year int) (models.BillDTO, error)
	GenerateMonthlyBills(ctx context.Context)
}
