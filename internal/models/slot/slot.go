package slot

import (
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type Slot struct {
	BuildingId  uuid.UUID
	FloorNumber int
	SlotNumber  int
	SlotType    vehicletypes.VehicleType
}

func (s Slot) GetID() string {
	return fmt.Sprintf("%v%v%v", s.BuildingId, s.FloorNumber, s.SlotNumber)
}
