package models

import (
	"fmt"
	"time"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
)

type ParkingHistoryDTO struct {
	TicketId     string
	NumberPlate  string
	BuildingId   string
	FLoorNumber  int
	SlotNumber   int
	StartTime    time.Time
	EndTime      time.Time
	VechicleType vehicletypes.VehicleType
}

func (phdto *ParkingHistoryDTO) String() string {
	return fmt.Sprintf("TicketId: %s\nNumberPlate: %s\nBuildingId: %s\nFloorNumber: %d\nSlotNumber: %d\nStartTime: %s\nEndTime: %s",
		phdto.TicketId, phdto.NumberPlate, phdto.BuildingId, phdto.FLoorNumber, phdto.SlotNumber, phdto.StartTime, phdto.EndTime)
}
