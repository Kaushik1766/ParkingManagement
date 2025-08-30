package slotrepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	slot "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSQLSlotRepository_AddSlot(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		buildingId  uuid.UUID
		floorNumber int
		slotNumber  int
		slotType    vehicletypes.VehicleType
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
				slotNumber:  1,
				slotType:    vehicletypes.FourWheeler,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "slots"`).
					WithArgs(buildingID, 1, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
				slotNumber:  1,
				slotType:    vehicletypes.FourWheeler,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "slots"`).
					WithArgs(buildingID, 1, 1, 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if err := tt.sqlsr.AddSlot(tt.args.buildingId, tt.args.floorNumber, tt.args.slotNumber, tt.args.slotType); (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.AddSlot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLSlotRepository_DeleteSlot(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		buildingId  uuid.UUID
		floorNumber int
		slotNumber  int
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
				slotNumber:  1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "slots"`).
					WithArgs(buildingID, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
				slotNumber:  1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "slots"`).
					WithArgs(buildingID, 1, 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if err := tt.sqlsr.DeleteSlot(tt.args.buildingId, tt.args.floorNumber, tt.args.slotNumber); (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.DeleteSlot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLSlotRepository_GetSlotsByFloor(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		buildingId  uuid.UUID
		floorNumber int
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		want      []slot.Slot
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"building_id", "floor_number", "slot_number", "slot_type"}).
					AddRow(buildingID, 1, 1, 1)
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 AND floor_number = \$2 ORDER BY slot_number asc`).
					WithArgs(buildingID, 1).
					WillReturnRows(rows)
			},
			want: []slot.Slot{
				{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 AND floor_number = \$2 ORDER BY slot_number asc`).
					WithArgs(buildingID, 1).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlsr.GetSlotsByFloor(tt.args.buildingId, tt.args.floorNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.GetSlotsByFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLSlotRepository.GetSlotsByFloor() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLSlotRepository_GetFreeSlotsByFloor(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		buildingId  uuid.UUID
		floorNumber int
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		want      []slot.Slot
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				slotRows := sqlmock.NewRows([]string{"building_id", "floor_number", "slot_number", "slot_type"}).
					AddRow(buildingID, 1, 1, 1)
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 AND floor_number = \$2 ORDER BY slot_number asc`).
					WithArgs(buildingID, 1).
					WillReturnRows(slotRows)
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE \("vehicles"."assigned_building_id","vehicles"."assigned_floor_number","vehicles"."assigned_slot_number"\) IN \(\(\$1,\$2,\$3\)\)`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "license_plate", "vehicle_type", "assigned_building_id", "assigned_floor_number", "assigned_slot_number"}))
			},
			want: []slot.Slot{
				{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.FourWheeler,
					Vehicles:    []slot.Vehicle{}, // Empty slice instead of nil
				},
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 AND floor_number = \$2 ORDER BY slot_number asc`).
					WithArgs(buildingID, 1).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlsr.GetFreeSlotsByFloor(tt.args.buildingId, tt.args.floorNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.GetFreeSlotsByFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLSlotRepository.GetFreeSlotsByFloor() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLSlotRepository_GetFreeSlotsByBuilding(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		buildingId uuid.UUID
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		want      []slot.Slot
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId: buildingID,
			},
			mockSetup: func() {
				slotRows := sqlmock.NewRows([]string{"building_id", "floor_number", "slot_number", "slot_type"}).
					AddRow(buildingID, 1, 1, 1).
					AddRow(buildingID, 1, 2, 0)
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 ORDER BY slot_number asc`).
					WithArgs(buildingID).
					WillReturnRows(slotRows)
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE \("vehicles"."assigned_building_id","vehicles"."assigned_floor_number","vehicles"."assigned_slot_number"\) IN \(\(\$1,\$2,\$3\),\(\$4,\$5,\$6\)\)`).
					WithArgs(buildingID, 1, 1, buildingID, 1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "license_plate", "vehicle_type", "assigned_building_id", "assigned_floor_number", "assigned_slot_number"}))
			},
			want: []slot.Slot{
				{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.FourWheeler,
					Vehicles:    []slot.Vehicle{}, // Empty slice instead of nil
				},
				{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  2,
					SlotType:    vehicletypes.TwoWheeler,
					Vehicles:    []slot.Vehicle{}, // Empty slice instead of nil
				},
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				buildingId: buildingID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE building_id = \$1 ORDER BY slot_number asc`).
					WithArgs(buildingID).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlsr.GetFreeSlotsByBuilding(tt.args.buildingId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.GetFreeSlotsByBuilding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLSlotRepository.GetFreeSlotsByBuilding() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLSlotRepository_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type args struct {
		slot slot.Slot
	}
	tests := []struct {
		name      string
		sqlsr     *SQLSlotRepository
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				slot: slot.Slot{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "slots" SET "slot_type"=\$1 WHERE "building_id" = \$2 AND "floor_number" = \$3 AND "slot_number" = \$4`).
					WithArgs(1, buildingID, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlsr: &SQLSlotRepository{
				db: gormDb,
			},
			args: args{
				slot: slot.Slot{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "slots" SET "slot_type"=\$1 WHERE "building_id" = \$2 AND "floor_number" = \$3 AND "slot_number" = \$4`).
					WithArgs(1, buildingID, 1, 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if err := tt.sqlsr.Save(tt.args.slot); (err != nil) != tt.wantErr {
				t.Errorf("SQLSlotRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNewSQLSlotRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *SQLSlotRepository
	}{
		{
			name: "success",
			args: args{
				db: gormDb,
			},
			want: &SQLSlotRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLSlotRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLSlotRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
