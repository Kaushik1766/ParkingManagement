package slotassignment

import (
	"context"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
)

func TestNewSlotAssignmentService(t *testing.T) {
	type args struct {
		vehicleRepo  vehiclerepository.VehicleStorage
		floorRepo    floorrepository.FloorStorage
		buildingRepo buildingrepository.BuildingStorage
		slotRepo     slotrepository.SlotStorage
		officeRepo   officerepository.OfficeStorage
	}
	tests := []struct {
		name string
		args args
		want *SlotAssignmentService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlotAssignmentService(tt.args.vehicleRepo, tt.args.floorRepo, tt.args.buildingRepo, tt.args.slotRepo, tt.args.officeRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlotAssignmentService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlotAssignmentService_AssignSlot(t *testing.T) {
	type fields struct {
		vehicleRepo  vehiclerepository.VehicleStorage
		floorRepo    floorrepository.FloorStorage
		buildingRepo buildingrepository.BuildingStorage
		slotRepo     slotrepository.SlotStorage
		officeRepo   officerepository.OfficeStorage
	}
	type args struct {
		ctx       context.Context
		vehicleId string
		slot      models.Slot
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sas := &SlotAssignmentService{
				vehicleRepo:  tt.fields.vehicleRepo,
				floorRepo:    tt.fields.floorRepo,
				buildingRepo: tt.fields.buildingRepo,
				slotRepo:     tt.fields.slotRepo,
				officeRepo:   tt.fields.officeRepo,
			}
			if err := sas.AssignSlot(tt.args.ctx, tt.args.vehicleId, tt.args.slot); (err != nil) != tt.wantErr {
				t.Errorf("AssignSlot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlotAssignmentService_AutoAssignSlot(t *testing.T) {
	type fields struct {
		vehicleRepo  vehiclerepository.VehicleStorage
		floorRepo    floorrepository.FloorStorage
		buildingRepo buildingrepository.BuildingStorage
		slotRepo     slotrepository.SlotStorage
		officeRepo   officerepository.OfficeStorage
	}
	type args struct {
		ctx       context.Context
		vehicleId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sas := &SlotAssignmentService{
				vehicleRepo:  tt.fields.vehicleRepo,
				floorRepo:    tt.fields.floorRepo,
				buildingRepo: tt.fields.buildingRepo,
				slotRepo:     tt.fields.slotRepo,
				officeRepo:   tt.fields.officeRepo,
			}
			if err := sas.AutoAssignSlot(tt.args.ctx, tt.args.vehicleId); (err != nil) != tt.wantErr {
				t.Errorf("AutoAssignSlot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlotAssignmentService_GetVehiclesWithUnassignedSlots(t *testing.T) {
	type fields struct {
		vehicleRepo  vehiclerepository.VehicleStorage
		floorRepo    floorrepository.FloorStorage
		buildingRepo buildingrepository.BuildingStorage
		slotRepo     slotrepository.SlotStorage
		officeRepo   officerepository.OfficeStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Vehicle
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sas := &SlotAssignmentService{
				vehicleRepo:  tt.fields.vehicleRepo,
				floorRepo:    tt.fields.floorRepo,
				buildingRepo: tt.fields.buildingRepo,
				slotRepo:     tt.fields.slotRepo,
				officeRepo:   tt.fields.officeRepo,
			}
			got, err := sas.GetVehiclesWithUnassignedSlots(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVehiclesWithUnassignedSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVehiclesWithUnassignedSlots() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlotAssignmentService_UnassignSlot(t *testing.T) {
	type fields struct {
		vehicleRepo  vehiclerepository.VehicleStorage
		floorRepo    floorrepository.FloorStorage
		buildingRepo buildingrepository.BuildingStorage
		slotRepo     slotrepository.SlotStorage
		officeRepo   officerepository.OfficeStorage
	}
	type args struct {
		ctx       context.Context
		vehicleId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sas := &SlotAssignmentService{
				vehicleRepo:  tt.fields.vehicleRepo,
				floorRepo:    tt.fields.floorRepo,
				buildingRepo: tt.fields.buildingRepo,
				slotRepo:     tt.fields.slotRepo,
				officeRepo:   tt.fields.officeRepo,
			}
			if err := sas.UnassignSlot(tt.args.ctx, tt.args.vehicleId); (err != nil) != tt.wantErr {
				t.Errorf("UnassignSlot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
