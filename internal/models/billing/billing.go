package billing

import (
	"time"

	"github.com/google/uuid"
)

type Billing struct {
	BillId      uuid.UUID
	UserId      uuid.UUID
	Month       int
	Year        int
	TotalAmount float32
	CreatedAt   time.Time
}

func (b Billing) GetID() string {
	return b.BillId.String()
}
