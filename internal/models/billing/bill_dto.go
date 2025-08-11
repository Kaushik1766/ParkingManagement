package billing

import (
	"fmt"

	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
)

type BillDTO struct {
	ParkingHistory []parkinghistory.ParkingHistoryDTO `json:"parking_history"`
	TotalAmount    float64                            `json:"total_amount"`
	BillDate       string                             `json:"bill_date"`
	UserId         string                             `json:"user_id"`
}

func (bdto *BillDTO) String() string {
	parkingHistoryStr := ""
	for _, val := range bdto.ParkingHistory {
		parkingHistoryStr += val.String() + "\n\n"
	}
	return fmt.Sprintf("\n\nParkingHistory:\n%v\n TotalAmount: %.2f\n BillDate: %s\n UserId: %s\n\n", parkingHistoryStr, bdto.TotalAmount, bdto.BillDate, bdto.UserId)
}
