package slotrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	slot "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLSlotRepository struct {
	db *gorm.DB
}

func (sqlsr *SQLSlotRepository) AddSlot(buildingId uuid.UUID, floorNumber int, slotNumber int, slotType vehicletypes.VehicleType) error {
	slot := slot.Slot{
		BuildingID:  buildingId,
		FloorNumber: floorNumber,
		SlotNumber:  slotNumber,
		SlotType:    slotType,
	}

	return sqlsr.db.Create(&slot).Error
}

func (sqlsr *SQLSlotRepository) DeleteSlot(buildingId uuid.UUID, floorNumber int, slotNumber int) error {
	return sqlsr.db.
		Where("building_id = ? AND floor_number = ? AND slot_number = ?", buildingId, floorNumber, slotNumber).
		Delete(&models.Slot{}).Error
}

func (sqlsr *SQLSlotRepository) GetSlotsByFloor(buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error) {
	var slots []slot.Slot
	err := sqlsr.db.
		Where("building_id = ? AND floor_number = ?", buildingId, floorNumber).
		Find(&slots).Error
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (sqlsr *SQLSlotRepository) GetFreeSlotsByFloor(buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error) {
	var slots []slot.Slot
	err := sqlsr.db.
		Where("building_id = ? AND floor_number = ?", buildingId, floorNumber).
		Preload("Vehicles").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}
	var res []slot.Slot
	for _, s := range slots {
		if len(s.Vehicles) == 0 {
			res = append(res, s)
		}
	}
	return res, nil
}

func (sqlsr *SQLSlotRepository) SetSlotOccupied(buildingId uuid.UUID, floorNumber int, slotNumber int, isOccupied bool) error {
	panic("not implemented") // TODO: Implement
}

func (sqlsr *SQLSlotRepository) GetFreeSlotsByBuilding(buildingId uuid.UUID) ([]slot.Slot, error) {
	var slots []slot.Slot
	err := sqlsr.db.
		Where("building_id = ? ", buildingId).
		Preload("Vehicles").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}

	var freeSlots []slot.Slot
	for _, s := range slots {
		if len(s.Vehicles) == 0 {
			freeSlots = append(freeSlots, s)
		}
	}
	return freeSlots, nil
}

func (sqlsr *SQLSlotRepository) Save(slot slot.Slot) error {
	return sqlsr.db.Save(&slot).Error
}

func NewSQLSlotRepository(db *gorm.DB) *SQLSlotRepository {
	return &SQLSlotRepository{
		db: db,
	}
}
