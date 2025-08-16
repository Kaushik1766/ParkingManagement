package buildingrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type FileBuildingRepository struct {
	*sync.Mutex
	buildings []models.Building
}

func (fbr *FileBuildingRepository) GetAllBuildings() ([]models.Building, error) {
	fbr.Lock()
	defer fbr.Unlock()
	// if len(fbr.buildings) == 0 {
	// 	return nil, fmt.Errorf("no buildings found")
	// }
	return fbr.buildings, nil
}

func (fbr *FileBuildingRepository) GetBuildingByID(buildingID uuid.UUID) (models.Building, error) {
	fbr.Lock()
	defer fbr.Unlock()
	for _, b := range fbr.buildings {
		if b.BuildingID == buildingID {
			return b, nil
		}
	}
	return models.Building{}, fmt.Errorf("building with id %s not found", buildingID)
}

func (fbr *FileBuildingRepository) GetBuildingByName(name string) (models.Building, error) {
	fbr.Lock()
	defer fbr.Unlock()
	for _, b := range fbr.buildings {
		if b.BuildingName == name {
			return b, nil
		}
	}
	return models.Building{}, fmt.Errorf("building with name %s not found", name)
}

func (fbr *FileBuildingRepository) AddBuilding(name string) error {
	fbr.Lock()
	defer fbr.Unlock()
	for _, b := range fbr.buildings {
		if b.BuildingName == name {
			return fmt.Errorf("building with name %s already exists", name)
		}
	}
	fbr.buildings = append(fbr.buildings, models.Building{
		BuildingName: name,
		BuildingID:   uuid.New(),
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
	data, err := os.ReadFile(config.BuildingsPath)
	if err != nil {
		os.WriteFile(config.BuildingsPath, []byte("[]"), 0666)
		data, err = json.Marshal([]models.Building{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var buildingData []models.Building
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
	err = os.WriteFile(config.BuildingsPath, data, 0666)
	return err
}
