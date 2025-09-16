package parkinghistoryrepository

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewSQLParkingRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *SQLParkingRepository
	}{
		{
			name: "success",
			db:   &gorm.DB{},
			want: &SQLParkingRepository{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSQLParkingRepository(tt.db)
			if got == nil || got.db == nil {
				t.Errorf("NewSQLParkingRepository() returned nil or nil db")
			}
		})
	}
}

func TestSQLParkingRepository_UnparkByNumberPlate(t *testing.T) {
	tests := []struct {
		name        string
		numberplate string
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
	}{
		{
			name:        "success",
			numberplate: "ABC123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "parking_histories" SET "end_time"=\$1 WHERE vehicles\.number_plate = \$2 AND end_time IS NULL`).
					WithArgs(sqlmock.AnyArg(), "ABC123").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:        "database_error",
			numberplate: "ABC123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "parking_histories" SET "end_time"=\$1 WHERE vehicles\.number_plate = \$2 AND end_time IS NULL`).
					WithArgs(sqlmock.AnyArg(), "ABC123").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			err := repo.UnparkByNumberPlate(tt.numberplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.UnparkByNumberPlate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLParkingRepository_AddParking(t *testing.T) {
	buildingID := uuid.New()
	floorNumber := 1
	slotNumber := 1

	tests := []struct {
		name      string
		vehicle   models.Vehicle
		mockSetup func(mock sqlmock.Sqlmock)
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			vehicle: models.Vehicle{
				VehicleID:           uuid.New(),
				NumberPlate:         "ABC123",
				VehicleType:         vehicletypes.FourWheeler,
				UserID:              uuid.New(),
				AssignedBuildingID:  buildingID,
				AssignedFloorNumber: floorNumber,
				AssignedSlotNumber:  slotNumber,
				AssignedSlot: models.Slot{
					BuildingID:  buildingID,
					FloorNumber: floorNumber,
					SlotNumber:  slotNumber,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on parking_histories\.vehicle_id = vehicles\.vehicle_id WHERE vehicles\.assigned_building_id = \$1 AND vehicles\.assigned_floor_number = \$2 AND vehicles\.assigned_slot_number = \$3 AND end_time IS NULL ORDER BY "parking_histories"\."parking_id" LIMIT \$4`).
					WithArgs(buildingID, floorNumber, slotNumber, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "slots"`).
					WithArgs(buildingID, floorNumber, slotNumber, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`INSERT INTO "vehicles"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "assigned_building_id", "assigned_floor_number", "assigned_slot_number"}).AddRow(uuid.New(), buildingID, floorNumber, slotNumber))
				mock.ExpectQuery(`INSERT INTO "parking_histories"`).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"parking_id", "start_time", "end_time"}).AddRow(uuid.New(), time.Now(), nil))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "vehicle_already_parked",
			vehicle: models.Vehicle{
				VehicleID:           uuid.New(),
				NumberPlate:         "ABC123",
				VehicleType:         vehicletypes.FourWheeler,
				UserID:              uuid.New(),
				AssignedBuildingID:  buildingID,
				AssignedFloorNumber: floorNumber,
				AssignedSlotNumber:  slotNumber,
				AssignedSlot: models.Slot{
					BuildingID:  buildingID,
					FloorNumber: floorNumber,
					SlotNumber:  slotNumber,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on parking_histories\.vehicle_id = vehicles\.vehicle_id WHERE vehicles\.assigned_building_id = \$1 AND vehicles\.assigned_floor_number = \$2 AND vehicles\.assigned_slot_number = \$3 AND end_time IS NULL ORDER BY "parking_histories"\."parking_id" LIMIT \$4`).
					WithArgs(buildingID, floorNumber, slotNumber, 1).
					WillReturnRows(sqlmock.NewRows([]string{"parking_id", "vehicle_id", "start_time", "end_time"}).AddRow(uuid.New(), uuid.New(), time.Now(), nil))
			},
			wantErr: true,
		},
		{
			name: "no_assigned_slot",
			vehicle: models.Vehicle{
				VehicleID:   uuid.New(),
				NumberPlate: "ABC123",
				VehicleType: vehicletypes.FourWheeler,
				UserID:      uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			got, err := repo.AddParking(tt.vehicle)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.AddParking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("SQLParkingRepository.AddParking() returned empty parking ID")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLParkingRepository_Unpark(t *testing.T) {
	parkingID := uuid.New().String()

	tests := []struct {
		name      string
		id        string
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			id:   parkingID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "parking_histories" SET "end_time"=\$1 WHERE parking_id = \$2 AND end_time IS NULL`).
					WithArgs(sqlmock.AnyArg(), parkingID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database_error",
			id:   parkingID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "parking_histories" SET "end_time"=\$1 WHERE parking_id = \$2 AND end_time IS NULL`).
					WithArgs(sqlmock.AnyArg(), parkingID).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			err := repo.Unpark(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.Unpark() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLParkingRepository_GetParkingHistoryByNumberPlate(t *testing.T) {
	buildingID := uuid.New()
	vehicleID := uuid.New()
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	tests := []struct {
		name        string
		numberplate string
		startTime   time.Time
		endTime     time.Time
		mockSetup   func(mock sqlmock.Sqlmock)
		want        []models.ParkingHistoryDTO
		wantErr     bool
	}{
		{
			name:        "success",
			numberplate: "ABC123",
			startTime:   startTime,
			endTime:     endTime,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"parking_id", "vehicle_id", "start_time", "end_time",
				}).AddRow(
					uuid.New(), vehicleID, startTime, endTime,
				)
				mock.ExpectQuery(`SELECT \* FROM "parking_histories" WHERE number_plate = \$1 AND start_time >= \$2 AND end_time <= \$3 AND end_time IS NOT NULL`).
					WithArgs("ABC123", startTime, endTime).
					WillReturnRows(rows)
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE "vehicles"\."vehicle_id" = \$1`).
					WithArgs(vehicleID).
					WillReturnRows(sqlmock.NewRows([]string{
						"vehicle_id", "number_plate", "vehicle_type", "user_id", "assigned_building_id", "assigned_floor_number", "assigned_slot_number",
					}).AddRow(
						vehicleID, "ABC123", 1, uuid.New(), buildingID, 1, 1,
					))
			},
			want: []models.ParkingHistoryDTO{
				{
					TicketId:     "",
					NumberPlate:  "ABC123",
					BuildingId:   buildingID.String(),
					FLoorNumber:  1,
					SlotNumber:   1,
					StartTime:    startTime.Local(),
					EndTime:      endTime.Local(),
					VechicleType: vehicletypes.FourWheeler,
				},
			},
			wantErr: false,
		},
		{
			name:        "database_error",
			numberplate: "ABC123",
			startTime:   startTime,
			endTime:     endTime,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "parking_histories" WHERE number_plate = \$1 AND start_time >= \$2 AND end_time <= \$3 AND end_time IS NOT NULL`).
					WithArgs("ABC123", startTime, endTime).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			got, err := repo.GetParkingHistoryByNumberPlate(tt.numberplate, tt.startTime, tt.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.GetParkingHistoryByNumberPlate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("SQLParkingRepository.GetParkingHistoryByNumberPlate() returned %d records, want %d", len(got), len(tt.want))
					return
				}
				if len(got) > 0 && len(tt.want) > 0 {
					if got[0].NumberPlate != tt.want[0].NumberPlate ||
						got[0].BuildingId != tt.want[0].BuildingId ||
						got[0].FLoorNumber != tt.want[0].FLoorNumber ||
						got[0].SlotNumber != tt.want[0].SlotNumber {
						t.Errorf("SQLParkingRepository.GetParkingHistoryByNumberPlate() = %v, want %v", got[0], tt.want[0])
					}
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLParkingRepository_GetParkingHistoryByUser(t *testing.T) {
	userID := uuid.New()
	buildingID := uuid.New()
	vehicleID := uuid.New()
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	tests := []struct {
		name      string
		userId    string
		startTime time.Time
		endTime   time.Time
		mockSetup func(mock sqlmock.Sqlmock)
		want      []models.ParkingHistoryDTO
		wantErr   bool
	}{
		{
			name:      "success",
			userId:    userID.String(),
			startTime: startTime,
			endTime:   endTime,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"parking_id", "vehicle_id", "start_time", "end_time",
				}).AddRow(
					uuid.New(), vehicleID, startTime, endTime,
				)
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on vehicles\.vehicle_id = parking_histories\.vehicle_id WHERE vehicles\.user_id = \$1 AND start_time >= \$2 AND end_time <= \$3 AND end_time IS NOT NULL`).
					WithArgs(userID, startTime, endTime).
					WillReturnRows(rows)
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE "vehicles"\."vehicle_id" = \$1`).
					WithArgs(vehicleID).
					WillReturnRows(sqlmock.NewRows([]string{
						"vehicle_id", "number_plate", "vehicle_type", "user_id", "assigned_building_id", "assigned_floor_number", "assigned_slot_number",
					}).AddRow(
						vehicleID, "ABC123", 1, userID, buildingID, 1, 1,
					))
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE \("slots"\."building_id","slots"\."floor_number","slots"\."slot_number"\) IN \(\(\$1,\$2,\$3\)\)`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"building_id", "floor_number", "slot_number", "slot_type",
					}).AddRow(
						buildingID, 1, 1, 1,
					))
			},
			want: []models.ParkingHistoryDTO{
				{
					TicketId:     "",
					NumberPlate:  "ABC123",
					BuildingId:   buildingID.String(),
					FLoorNumber:  1,
					SlotNumber:   1,
					StartTime:    startTime.Local(),
					EndTime:      endTime.Local(),
					VechicleType: vehicletypes.FourWheeler,
				},
			},
			wantErr: false,
		},
		{
			name:      "database_error",
			userId:    userID.String(),
			startTime: startTime,
			endTime:   endTime,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on vehicles\.vehicle_id = parking_histories\.vehicle_id WHERE vehicles\.user_id = \$1 AND start_time >= \$2 AND end_time <= \$3 AND end_time IS NOT NULL`).
					WithArgs(userID, startTime, endTime).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			got, err := repo.GetParkingHistoryByUser(tt.userId, tt.startTime, tt.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.GetParkingHistoryByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("SQLParkingRepository.GetParkingHistoryByUser() returned %d records, want %d", len(got), len(tt.want))
					return
				}
				if len(got) > 0 && len(tt.want) > 0 {
					if got[0].NumberPlate != tt.want[0].NumberPlate ||
						got[0].BuildingId != tt.want[0].BuildingId ||
						got[0].FLoorNumber != tt.want[0].FLoorNumber ||
						got[0].SlotNumber != tt.want[0].SlotNumber {
						t.Errorf("SQLParkingRepository.GetParkingHistoryByUser() = %v, want %v", got[0], tt.want[0])
					}
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLParkingRepository_GetActiveUserParkings(t *testing.T) {
	userID := uuid.New()
	buildingID := uuid.New()
	vehicleID := uuid.New()
	startTime := time.Now().Add(-2 * time.Hour)

	tests := []struct {
		name      string
		userId    string
		mockSetup func(mock sqlmock.Sqlmock)
		want      []models.ParkingHistoryDTO
		wantErr   bool
	}{
		{
			name:   "success",
			userId: userID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"parking_id", "vehicle_id", "start_time", "end_time",
				}).AddRow(
					uuid.New(), vehicleID, startTime, nil,
				)
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on vehicles\.vehicle_id = parking_histories\.vehicle_id WHERE vehicles\.user_id = \$1 AND end_time IS NULL`).
					WithArgs(userID).
					WillReturnRows(rows)
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE "vehicles"\."vehicle_id" = \$1`).
					WithArgs(vehicleID).
					WillReturnRows(sqlmock.NewRows([]string{
						"vehicle_id", "number_plate", "vehicle_type", "user_id", "assigned_building_id", "assigned_floor_number", "assigned_slot_number",
					}).AddRow(
						vehicleID, "ABC123", 1, userID, buildingID, 1, 1,
					))
				mock.ExpectQuery(`SELECT \* FROM "slots" WHERE \("slots"\."building_id","slots"\."floor_number","slots"\."slot_number"\) IN \(\(\$1,\$2,\$3\)\)`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"building_id", "floor_number", "slot_number", "slot_type",
					}).AddRow(
						buildingID, 1, 1, 1,
					))
			},
			want: []models.ParkingHistoryDTO{
				{
					TicketId:     "",
					NumberPlate:  "ABC123",
					BuildingId:   buildingID.String() + "_1_1",
					FLoorNumber:  1,
					SlotNumber:   1,
					StartTime:    startTime.Local(),
					VechicleType: vehicletypes.FourWheeler,
				},
			},
			wantErr: false,
		},
		{
			name:   "database_error",
			userId: userID.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "parking_histories"\."parking_id","parking_histories"\."vehicle_id","parking_histories"\."start_time","parking_histories"\."end_time" FROM "parking_histories" left join vehicles on vehicles\.vehicle_id = parking_histories\.vehicle_id WHERE vehicles\.user_id = \$1 AND end_time IS NULL`).
					WithArgs(userID).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gormDb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := NewSQLParkingRepository(gormDb)

			tt.mockSetup(mock)
			got, err := repo.GetActiveUserParkings(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLParkingRepository.GetActiveUserParkings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("SQLParkingRepository.GetActiveUserParkings() returned %d records, want %d", len(got), len(tt.want))
					return
				}
				if len(got) > 0 && len(tt.want) > 0 {
					if got[0].NumberPlate != tt.want[0].NumberPlate ||
						got[0].BuildingId != tt.want[0].BuildingId ||
						got[0].FLoorNumber != tt.want[0].FLoorNumber ||
						got[0].SlotNumber != tt.want[0].SlotNumber {
						t.Errorf("SQLParkingRepository.GetActiveUserParkings() = %v, want %v", got[0], tt.want[0])
					}
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
