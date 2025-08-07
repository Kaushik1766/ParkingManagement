package vehiclerepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/google/uuid"
)

type FileVehicleRepository struct {
	*sync.Mutex
	vehicles []vehicle.Vehicle
	userRepo userrepository.UserStorage
}

func (fvr *FileVehicleRepository) Save(vehicle vehicle.Vehicle) error {
	fvr.Lock()
	defer fvr.Unlock()
	for i, val := range fvr.vehicles {
		if val.VehicleId == vehicle.VehicleId {
			fvr.vehicles[i] = vehicle
			return nil
		}
	}
	fvr.vehicles = append(fvr.vehicles, vehicle)
	return nil
}

func (fvr *FileVehicleRepository) AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (vehicle.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	for _, val := range fvr.vehicles {
		if val.NumberPlate == numberplate {
			return vehicle.Vehicle{}, errors.New("numberplate already registered")
		}
	}
	// TODO: add assigned slot and update interface
	fvr.vehicles = append(fvr.vehicles, vehicle.Vehicle{
		VehicleId:   uuid.New(),
		NumberPlate: numberplate,
		VehicleType: vehicleType,
		UserId:      userid,
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

func (fvr *FileVehicleRepository) GetVehicleById(vehicleId uuid.UUID) (vehicle.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	for _, val := range fvr.vehicles {
		if val.VehicleId == vehicleId {
			return val, nil
		}
	}
	return vehicle.Vehicle{}, errors.New("vehicle not found")
}

func (fvr *FileVehicleRepository) GetVehiclesByUserId(userId uuid.UUID) ([]vehicle.Vehicle, error) {
	fvr.Lock()
	defer fvr.Unlock()
	var result []vehicle.Vehicle
	// fmt.Println(fvr.vehicles)
	// fmt.Println(userId.String())
	for _, val := range fvr.vehicles {
		if val.UserId == userId {
			result = append(result, val)
		}
	}
	return result, nil
}

func NewFileVehicleRepository(userRepo userrepository.UserStorage) *FileVehicleRepository {
	data, err := os.ReadFile("vehicles.json")
	if err != nil {
		os.WriteFile("vehicles.json", []byte("[]"), 0666)
		data, err = json.Marshal([]vehicle.Vehicle{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var vehicleData []vehicle.Vehicle

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
	err = os.WriteFile("vehicles.json", data, 0666)
	return err
}
