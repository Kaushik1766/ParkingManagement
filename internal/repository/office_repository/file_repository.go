package officerepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type FileOfficeRepository struct {
	*sync.Mutex
	offices []models.Office
}

func (fr *FileOfficeRepository) GetBuildingAndFloorByOffice(officeName string) (uuid.UUID, int, error) {
	fr.Lock()
	defer fr.Unlock()
	for _, val := range fr.offices {
		if val.OfficeName == officeName {
			return val.BuildingID, val.FloorNumber, nil
		}
	}
	return uuid.Nil, 0, errors.New("office not found")
}

func NewFileOfficeRepository() *FileOfficeRepository {
	data, err := os.ReadFile(config.OfficesPath)
	if err != nil {
		os.WriteFile(config.OfficesPath, []byte("[]"), 0666)
		data, err = json.Marshal([]models.Office{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var officeData []models.Office
	err = json.Unmarshal(data, &officeData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileOfficeRepository{
		Mutex:   &sync.Mutex{},
		offices: officeData,
	}
}

func (fr *FileOfficeRepository) AddOffice(officeName string, buildingID uuid.UUID, floorNumber int) error {
	fr.Lock()
	defer fr.Unlock()
	for _, val := range fr.offices {
		if val.OfficeName == officeName {
			return errors.New("office already exists")
		}
	}
	fr.offices = append(fr.offices, models.Office{
		OfficeID:    uuid.New(),
		OfficeName:  officeName,
		BuildingID:  buildingID,
		FloorNumber: floorNumber,
	})
	return nil
}

func (fr *FileOfficeRepository) DeleteOffice(officeName string) error {
	fr.Lock()
	defer fr.Unlock()
	for i, val := range fr.offices {
		if val.OfficeName == officeName {
			fr.offices = append(fr.offices[:i], fr.offices[i+1:]...)
			return nil
		}
	}
	return errors.New("office not found")
}

func (fr *FileOfficeRepository) GetOfficesByBuilding(buildingID uuid.UUID) ([]models.Office, error) {
	fr.Lock()
	defer fr.Unlock()
	var offices []models.Office
	for _, val := range fr.offices {
		if val.BuildingID == buildingID {
			offices = append(offices, val)
		}
	}
	return offices, nil
}

func (fr *FileOfficeRepository) GetAllOffices() ([]models.Office, error) {
	fr.Lock()
	defer fr.Unlock()
	return fr.offices, nil
}

func (fr *FileOfficeRepository) GetOfficeByName(officeName string) (models.Office, error) {
	fr.Lock()
	defer fr.Unlock()
	for _, val := range fr.offices {
		if val.OfficeName == officeName {
			return val, nil
		}
	}
	return models.Office{}, errors.New("office not found")
}

func (fr *FileOfficeRepository) SerializeData() error {
	fr.Lock()
	defer fr.Unlock()
	data, err := json.Marshal(fr.offices)
	if err != nil {
		return err
	}
	err = os.WriteFile(config.OfficesPath, data, 0666)
	return err
}
