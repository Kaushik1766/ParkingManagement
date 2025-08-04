package officerepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/models/office"
)

type FileOfficeRepository struct {
	*sync.Mutex
	offices []office.Office
}

func (fr *FileOfficeRepository) GetBuildingAndFloorByOffice(officeName string) (string, int, error) {
	fr.Lock()
	defer fr.Unlock()
	for _, val := range fr.offices {
		if val.OfficeName == officeName {
			return val.BuildingName, val.FloorNumber, nil
		}
	}
	return "", 0, errors.New("office not found")
}

func NewFileOfficeRepository() *FileOfficeRepository {
	data, err := os.ReadFile("offices.json")
	if err != nil {
		os.WriteFile("offices.json", []byte("[]"), 0666)
		data, err = json.Marshal([]office.Office{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var officeData []office.Office
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

func (fr *FileOfficeRepository) AddOffice(officeName string, buildingName string, floorNumber int) error {
	fr.Lock()
	defer fr.Unlock()
	for _, val := range fr.offices {
		if val.OfficeName == officeName {
			return errors.New("office already exists")
		}
	}
	fr.offices = append(fr.offices, office.Office{
		OfficeName:   officeName,
		BuildingName: buildingName,
		FloorNumber:  floorNumber,
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

func (fr *FileOfficeRepository) SerializeData() error {
	fr.Lock()
	defer fr.Unlock()
	data, err := json.Marshal(fr.offices)
	if err != nil {
		return err
	}
	err = os.WriteFile("offices.json", data, 0666)
	return err
}
