package officeservice

import (
	"context"
	"errors"

	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
)

type OfficeService struct {
	officeRepo   officerepository.OfficeStorage
	buildingRepo buildingrepository.BuildingStorage
	flooRepo     floorrepository.FloorStorage
}

func NewOfficeService(officeRepo officerepository.OfficeStorage,
	buildingRepo buildingrepository.BuildingStorage,
	flooRepo floorrepository.FloorStorage,
) *OfficeService {
	return &OfficeService{
		officeRepo:   officeRepo,
		buildingRepo: buildingRepo,
		flooRepo:     flooRepo,
	}
}

func (officeServ *OfficeService) AddOffice(ctx context.Context, officeName string, buildingName string, floorNumber int) error {
	if officeName == "" || buildingName == "" || floorNumber <= 0 {
		return errors.New("invalid input parameters")
	}

	building, err := officeServ.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return errors.New("building does not exist")
	}

	_, err = officeServ.flooRepo.GetFloor(building.BuildingId, floorNumber)
	if err != nil {
		return errors.New("floor does not exist in the specified building")
	}

	return officeServ.officeRepo.AddOffice(officeName, buildingName, floorNumber)
}

func (officeServ *OfficeService) RemoveOffice(ctx context.Context, officeName string) error {
	return officeServ.officeRepo.DeleteOffice(officeName)
}

func (officeServ *OfficeService) ListOfficesByBuilding(ctx context.Context, buildingName string) (map[int][]string, error) {
	offices, err := officeServ.officeRepo.GetOfficesByBuilding(buildingName)
	if err != nil {
		return nil, errors.New("no offices in building")
	}

	officeMap := make(map[int][]string)
	for _, office := range offices {
		officeMap[office.FloorNumber] = append(officeMap[office.FloorNumber], office.OfficeName)
	}
	return officeMap, nil
}

func (officeServ *OfficeService) GetAllOfficeNames(ctx context.Context) ([]string, error) {
	offices, err := officeServ.officeRepo.GetAllOffices()
	if err != nil {
		return nil, errors.New("no offices found")
	}

	var officeNames []string
	for _, office := range offices {
		officeNames = append(officeNames, office.OfficeName)
	}
	return officeNames, nil
}
