package floorrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/models/floor"
	"github.com/google/uuid"
)

type FileFloorRepository struct {
	*sync.Mutex
	floors []floor.Floor
}

func (ffr *FileFloorRepository) GetFloorsByBuildingId(buildingId uuid.UUID) ([]int, error) {
	ffr.Lock()
	defer ffr.Unlock()
	var floorNumbers []int
	for _, f := range ffr.floors {
		if f.BuildingId == buildingId {
			floorNumbers = append(floorNumbers, f.FloorNumber)
		}
	}
	if len(floorNumbers) == 0 {
		return nil, fmt.Errorf("no floors found for building %s", buildingId)
	}
	return floorNumbers, nil
}

func (ffr *FileFloorRepository) GetFloor(buildingId uuid.UUID, floorNumber int) (int, error) {
	ffr.Lock()
	defer ffr.Unlock()
	for _, f := range ffr.floors {
		if f.BuildingId == buildingId && f.FloorNumber == floorNumber {
			return f.FloorNumber, nil
		}
	}
	return 0, fmt.Errorf("floor %d for building %s not found", floorNumber, buildingId)
}

func (ffr *FileFloorRepository) AddFloor(buildingId uuid.UUID, floorNumber int) error {
	ffr.Lock()
	defer ffr.Unlock()
	for _, f := range ffr.floors {
		if f.BuildingId == buildingId && f.FloorNumber == floorNumber {
			return fmt.Errorf("floor %d already exists", floorNumber)
		}
	}
	ffr.floors = append(ffr.floors, floor.Floor{
		BuildingId:  buildingId,
		FloorNumber: floorNumber,
	})
	return nil
}

func (ffr *FileFloorRepository) DeleteFloor(buildingId uuid.UUID, floorNumber int) error {
	ffr.Lock()
	defer ffr.Unlock()
	for i, f := range ffr.floors {
		if f.BuildingId == buildingId && f.FloorNumber == floorNumber {
			ffr.floors = append(ffr.floors[:i], ffr.floors[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("floor %d for building %s not found", floorNumber, buildingId)
}

func NewFileFloorRepository() *FileFloorRepository {
	data, err := os.ReadFile("floors.json")
	if err != nil {
		os.WriteFile("floors.json", []byte("[]"), 0666)
		data, err = json.Marshal([]floor.Floor{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var floorData []floor.Floor
	err = json.Unmarshal(data, &floorData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileFloorRepository{
		Mutex:  &sync.Mutex{},
		floors: floorData,
	}
}

func (ffr *FileFloorRepository) SerializeData() {
	ffr.Lock()
	defer ffr.Unlock()
	data, err := json.Marshal(ffr.floors)
	if err != nil {
		fmt.Println("unable to marshal floors data")
		return
	}
	err = os.WriteFile("floors.json", data, 0666)
	if err != nil {
		fmt.Println("unable to write floors data to file")
	}
}
