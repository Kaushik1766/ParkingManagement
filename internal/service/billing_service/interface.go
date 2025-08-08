package billingservice

type BillingMgr interface {
	GenerateMonthlyInvoice(customerID string) (string, error)
}
