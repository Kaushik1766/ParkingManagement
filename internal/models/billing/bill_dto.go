package billing

import (
	"fmt"

	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
)

type BillDTO struct {
	ParkingHistory []parkinghistory.ParkingHistoryDTO `json:"parking_history"`
	TotalAmount    float64                            `json:"total_amount"`
	BillDate       string                             `json:"bill_date"`
}

func (bdto *BillDTO) String() string {
	var parkingHistoryStr string
	for _, ph := range bdto.ParkingHistory {
		parkingHistoryStr += ph.String() + "\n"
	}
	return "BillDTO{" +
		"ParkingHistory: [" + parkingHistoryStr + "], " +
		"TotalAmount: " + fmt.Sprintf("%.2f", bdto.TotalAmount) + ", " +
		"BillDate: " + bdto.BillDate
}
