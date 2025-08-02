package slot

import (
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/google/uuid"
)

type Slot struct {
	BuildingId  uuid.UUID
	FloorNumber int
	SlotNumber  int
	SlotType    enums.VehicleType
}

func (s Slot) GetID() string {
	return fmt.Sprintf("%v%v%v", s.BuildingId, s.FloorNumber, s.SlotNumber)
}
