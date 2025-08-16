package vehiclerepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/google/uuid"
)

type FileVehicleRepository struct {
	*sync.Mutex
	vehicles []models.Vehicle
	userRepo userrepository.UserStorage
}

func (fvr *FileVehicleRepository) Save(vehicle models.Vehicle) error {
	fvr.Lock()
	defer fvr.Unlock()
	for i, val := range fvr.vehicles {
		if val.VehicleID == vehicle.VehicleID {
			fvr.vehicles[i] = vehicle
			return nil
		}
	}
	fvr.vehicles = append(fvr.vehicles, vehicle)
	return nil
}

func (fvr *FileVehicleRepository) GetVehicleByNumberPlate(numberplate string) (models.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	// fmt.Println(fvr.vehicles)
	// fmt.Println(numberplate)
	for _, val := range fvr.vehicles {
		if val.NumberPlate == numberplate && val.IsActive {
			return val, nil
		}
	}
	return models.Vehicle{}, errors.New("vehicle not found")
}

func (fvr *FileVehicleRepository) GetVehiclesWithUnassignedSlots() ([]models.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	var result []models.Vehicle
	for _, val := range fvr.vehicles {
		if val.IsActive && val.AssignedSlot.BuildingID == uuid.Nil {
			result = append(result, val)
		}
	}
	return result, nil
}

func (fvr *FileVehicleRepository) AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	for _, val := range fvr.vehicles {
		if val.NumberPlate == numberplate {
			return models.Vehicle{}, errors.New("numberplate already registered")
		}
	}
	// TODO: add assigned slot and update interface
	fvr.vehicles = append(fvr.vehicles, models.Vehicle{
		VehicleID:   uuid.New(),
		NumberPlate: numberplate,
		VehicleType: vehicleType,
		UserID:      userid,
		IsActive:    true,
	})
	return fvr.vehicles[len(fvr.vehicles)-1], nil
}

func (fvr *FileVehicleRepository) RemoveVehicle(numberplate string) error {
	fvr.Lock()
	defer fvr.Unlock()
	for i, val := range fvr.vehicles {
		if val.NumberPlate == numberplate {
			// fvr.vehicles = append(fvr.vehicles[:i], fvr.vehicles[i+1:]...)
			fvr.vehicles[i].IsActive = false
			return nil
		}
	}
	return errors.New("numberplate not found")
}

func (fvr *FileVehicleRepository) GetVehicleById(vehicleId uuid.UUID) (models.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	for _, val := range fvr.vehicles {
		if val.VehicleID == vehicleId {
			return val, nil
		}
	}
	return models.Vehicle{}, errors.New("vehicle not found")
}

func (fvr *FileVehicleRepository) GetVehiclesByUserId(userId uuid.UUID) ([]models.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	var result []models.Vehicle
	// fmt.Println(fvr.vehicles)
	// fmt.Println(userId.String())
	for _, val := range fvr.vehicles {
		if val.UserID == userId && val.IsActive {
			result = append(result, val)
		}
	}
	return result, nil
}

func NewFileVehicleRepository(userRepo userrepository.UserStorage) *FileVehicleRepository {
	data, err := os.ReadFile(config.VehiclesPath)
	if err != nil {
		os.WriteFile(config.VehiclesPath, []byte("[]"), 0666)
		data, err = json.Marshal([]models.Vehicle{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var vehicleData []models.Vehicle

	err = json.Unmarshal(data, &vehicleData)
	if err != nil {
		panic("corrupted data")
	}
	return &FileVehicleRepository{
		Mutex:    &sync.Mutex{},
		vehicles: vehicleData,
		userRepo: userRepo,
	}
}

func (fvr *FileVehicleRepository) SerializeData() error {
	fvr.Lock()
	defer fvr.Unlock()
	data, err := json.Marshal(fvr.vehicles)
	if err != nil {
		return err
	}
	err = os.WriteFile(config.VehiclesPath, data, 0666)
	return err
}
