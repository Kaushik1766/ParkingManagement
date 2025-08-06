package slotrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type FileSlotRepository struct {
	*sync.Mutex
	slots []slot.Slot
}

func (fsr *FileSlotRepository) GetSlotsByFloor(buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error) {
	fsr.Lock()
	defer fsr.Unlock()
	var slots []slot.Slot
	for _, s := range fsr.slots {
		if s.BuildingId == buildingId && s.FloorNumber == floorNumber {
			slots = append(slots, s)
		}
	}
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].SlotNumber < slots[j].SlotNumber
	})
	// if len(slots) == 0 {
	// 	return nil, fmt.Errorf("no slots found for building %s, floor %d", buildingId, floorNumber)
	// }
	return slots, nil
}

func (fsr *FileSlotRepository) AddSlot(buildingId uuid.UUID, floorNumber int, slotNumber int, slotType vehicletypes.VehicleType) error {
	fsr.Lock()
	defer fsr.Unlock()
	for _, s := range fsr.slots {
		if s.BuildingId == buildingId && s.FloorNumber == floorNumber && s.SlotNumber == slotNumber {
			return fmt.Errorf("slot already exists at building %s, floor %d, slot %d", buildingId, floorNumber, slotNumber)
		}
	}
	fsr.slots = append(fsr.slots, slot.Slot{
		BuildingId:  buildingId,
		FloorNumber: floorNumber,
		SlotNumber:  slotNumber,
		SlotType:    slotType,
		IsOccupied:  false,
	})
	return nil
}

func (fsr *FileSlotRepository) DeleteSlot(buildingId uuid.UUID, floorNumber int, slotNumber int) error {
	fsr.Lock()
	defer fsr.Unlock()
	for i, s := range fsr.slots {
		if s.BuildingId == buildingId && s.FloorNumber == floorNumber && s.SlotNumber == slotNumber {
			fsr.slots = append(fsr.slots[:i], fsr.slots[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("slot not found at building %s, floor %d, slot %d", buildingId, floorNumber, slotNumber)
}

func NewFileSlotRepository() *FileSlotRepository {
	data, err := os.ReadFile("slots.json")
	if err != nil {
		os.WriteFile("slots.json", []byte("[]"), 0666)
		data, err = json.Marshal([]slot.Slot{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var slotData []slot.Slot
	err = json.Unmarshal(data, &slotData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileSlotRepository{
		Mutex: &sync.Mutex{},
		slots: slotData,
	}
}

func (fsr *FileSlotRepository) SerializeData() {
	data, err := json.Marshal(fsr.slots)
	if err != nil {
		fmt.Println("unable to marshal slot data")
		return
	}
	err = os.WriteFile("slots.json", data, 0666)
	if err != nil {
		fmt.Println("unable to write slot data to file")
	}
}
