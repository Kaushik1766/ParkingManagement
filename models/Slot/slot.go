package slot

import (
	"github.com/Kaushik1766/ParkingManagement/models/enums"
	"github.com/google/uuid"
)

type Slot struct {
	BuildingId  uuid.UUID
	FloorNumber int
	SlotNumber  int
	SlotType    enums.VehicleType
}
