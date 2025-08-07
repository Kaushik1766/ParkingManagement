package slot

import (
	"fmt"
	"strconv"
	"strings"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

// Slot is slot, isoccupied is for is assigned
type Slot struct {
	BuildingId  uuid.UUID
	FloorNumber int
	SlotNumber  int
	SlotType    vehicletypes.VehicleType
	IsOccupied  bool
}

func (s Slot) GetID() string {
	return fmt.Sprintf("%v%v%v", s.BuildingId, s.FloorNumber, s.SlotNumber)
}

func (s Slot) String() string {
	if s.BuildingId == uuid.Nil {
		return "unassigned"
	}
	return fmt.Sprintf("%v_%v_%v", s.BuildingId, s.FloorNumber, s.SlotNumber)
}

func (s Slot) ToIdentifiableSlot(slotString string) (*Slot, error) {
	parts := strings.Split(slotString, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid slot string format: %s", slotString)
	}
	buildingId, err := uuid.Parse(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid building ID: %s", parts[0])
	}
	floorNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid floor number: %s", parts[1])
	}
	slotNumber, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid slot number: %s", parts[2])
	}
	return &Slot{
		BuildingId:  buildingId,
		FloorNumber: floorNumber,
		SlotNumber:  slotNumber,
	}, nil
}
