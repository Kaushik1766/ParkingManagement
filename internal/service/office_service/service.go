package officeservice

import (
	"context"
	"errors"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	"github.com/google/uuid"
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

func (officeServ *OfficeService) AddOffice(ctx context.Context, officeName string, buildingId string, floorNumber int) error {
	if officeName == "" || buildingId == "" || floorNumber <= 0 {
		return errors.New("invalid input parameters")
	}

	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return err
	}

	_, err = officeServ.flooRepo.GetFloor(buildingUUID, floorNumber)
	if err != nil {
		return errors.New("floor does not exist in the specified building")
	}

	return officeServ.officeRepo.AddOffice(officeName, buildingUUID, floorNumber)
}

func (officeServ *OfficeService) RemoveOffice(ctx context.Context, officeId string) error {
	return officeServ.officeRepo.DeleteOffice(officeId)
}

func (officeServ *OfficeService) ListOfficesByBuilding(ctx context.Context, buildingId string) ([]models.OfficeDTO, error) {
	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return nil, errors.New("invalid building ID format")
	}

	offices, err := officeServ.officeRepo.GetOfficesByBuilding(buildingUUID)
	if err != nil {
		return nil, errors.New("no offices in building")
	}

	var officeDTOs []models.OfficeDTO
	for _, office := range offices {
		officeDTO := models.OfficeDTO{
			OfficeName:  office.OfficeName,
			BuildingID:  office.BuildingID.String(),
			FloorNumber: office.FloorNumber,
			OfficeID:    office.OfficeID.String(),
		}
		officeDTOs = append(officeDTOs, officeDTO)
	}

	return officeDTOs, nil
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

func (officeServ *OfficeService) GetOfficeByName(ctx context.Context, officeName string) (models.Office, error) {
	if officeName == "" {
		return models.Office{}, errors.New("office name cannot be empty")
	}

	officeStruct, err := officeServ.officeRepo.GetOfficeByName(officeName)
	if err != nil {
		return models.Office{}, errors.New("office does not exist")
	}
	return officeStruct, nil
}
