package parkinghistory

import (
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
)

type ParkingHistoryDTO struct {
	TicketId     string
	NumberPlate  string
	BuildingId   string
	FLoorNumber  int
	SlotNumber   int
	StartTime    string
	EndTime      string
	VechicleType vehicletypes.VehicleType
}

func (phdto *ParkingHistoryDTO) String() string {
	return fmt.Sprintf("TicketId: %s, NumberPlate: %s, BuildingId: %s, FloorNumber: %d, SlotNumber: %d, StartTime: %s, EndTime: %s",
		phdto.TicketId, phdto.NumberPlate, phdto.BuildingId, phdto.FLoorNumber, phdto.SlotNumber, phdto.StartTime, phdto.EndTime)
}
