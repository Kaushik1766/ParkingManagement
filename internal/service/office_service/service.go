package officeservice

import (
	"context"
	"errors"
	"log"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
)

type OfficeService struct {
	officeRepo officerepository.OfficeStorage
}

func NewOfficeService(officeRepo officerepository.OfficeStorage,
) *OfficeService {
	return &OfficeService{
		officeRepo: officeRepo,
	}
}

func (officeServ *OfficeService) AddOffice(ctx context.Context, officeName string, buildingId string, floorNumber int) error {
	if officeName == "" || buildingId == "" || floorNumber <= 0 {
		return errors.New("invalid input parameters")
	}

	return officeServ.officeRepo.AddOffice(ctx, officeName, buildingId, floorNumber)
}

func (officeServ *OfficeService) RemoveOffice(ctx context.Context, officeId string) error {
	return officeServ.officeRepo.DeleteOffice(ctx, officeId)
}

func (officeServ *OfficeService) ListOfficesByBuilding(ctx context.Context, buildingId string) ([]models.OfficeDTO, error) {
	offices, err := officeServ.officeRepo.GetOfficesByBuilding(ctx, buildingId)
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

func (officeServ *OfficeService) GetAllOffices(ctx context.Context) ([]models.OfficeDTO, error) {
	offices, err := officeServ.officeRepo.GetAllOffices(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("no offices found")
	}

	offceDtOs := make([]models.OfficeDTO, 0, len(offices))
	for _, office := range offices {
		officeDTO := models.OfficeDTO{
			OfficeName:  office.OfficeName,
			BuildingID:  office.BuildingID.String(),
			FloorNumber: office.FloorNumber,
			OfficeID:    office.OfficeID.String(),
		}
		offceDtOs = append(offceDtOs, officeDTO)
	}

	return offceDtOs, nil
}
