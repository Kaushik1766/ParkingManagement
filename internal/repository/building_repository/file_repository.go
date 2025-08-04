package buildingrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/models/building"
	"github.com/google/uuid"
)

type FileBuildingRepository struct {
	*sync.Mutex
	buildings []building.Building
}

func (fbr *FileBuildingRepository) GetBuildingByName(name string) (building.Building, error) {
	fbr.Lock()
	defer fbr.Unlock()
	for _, b := range fbr.buildings {
		if b.BuildingName == name {
			return b, nil
		}
	}
	return building.Building{}, fmt.Errorf("building with name %s not found", name)
}

func (fbr *FileBuildingRepository) AddBuilding(name string) error {
	fbr.Lock()
	defer fbr.Unlock()
	for _, b := range fbr.buildings {
		if b.BuildingName == name {
			return fmt.Errorf("building with name %s already exists", name)
		}
	}
	fbr.buildings = append(fbr.buildings, building.Building{
		BuildingName: name,
		BuildingId:   uuid.New(),
	})
	return nil
}

func (fbr *FileBuildingRepository) DeleteBuilding(name string) error {
	fbr.Lock()
	defer fbr.Unlock()
	for i, b := range fbr.buildings {
		if b.BuildingName == name {
			fbr.buildings = append(fbr.buildings[:i], fbr.buildings[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("building with name %s not found", name)
}

func NewFileBuildingRepository() *FileBuildingRepository {
	data, err := os.ReadFile("buildings.json")
	if err != nil {
		os.WriteFile("buildings.json", []byte("[]"), 0666)
		data, err = json.Marshal([]building.Building{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var buildingData []building.Building
	err = json.Unmarshal(data, &buildingData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileBuildingRepository{
		Mutex:     &sync.Mutex{},
		buildings: buildingData,
	}
}

func (fbr *FileBuildingRepository) SerializeData() error {
	fbr.Lock()
	defer fbr.Unlock()
	data, err := json.Marshal(fbr.buildings)
	if err != nil {
		return err
	}
	err = os.WriteFile("buildings.json", data, 0666)
	return err
}
