package billingservice

import "context"

//go:generate mockgen -source=interface.go -destination=../../../mocks/billing_service_mock.go -package=mocks
type BillingMgr interface {
	GenerateMonthlyInvoice(ctx context.Context)
}
