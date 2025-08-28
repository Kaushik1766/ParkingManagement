package officerepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewSQLOfficeRepository(t *testing.T) {
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
		want *SQLOfficeRepository
	}{
		{
			name: "success",
			args: args{
				db: gormDb,
			},
			want: &SQLOfficeRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLOfficeRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLOfficeRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLOfficeRepository_AddOffice(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()
	officeID := uuid.New()

	type args struct {
		officeName  string
		buildingID  string
		floorNumber int
	}
	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName:  "wg",
				buildingID:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "offices"`).
					WithArgs("wg", buildingID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id"}).
						AddRow(officeID))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "invalid_building_id",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName:  "wg",
				buildingID:  "invalid-uuid",
				floorNumber: 1,
			},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName:  "wg",
				buildingID:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "offices"`).
					WithArgs("wg", buildingID, 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if err := tt.sqlor.AddOffice(tt.args.officeName, tt.args.buildingID, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.AddOffice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLOfficeRepository_DeleteOffice(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	officeID := uuid.New()

	type args struct {
		officeId string
	}
	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeId: officeID.String(),
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "offices"`).
					WithArgs(officeID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "invalid_office_id",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeId: "invalid-uuid",
			},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeId: officeID.String(),
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "offices"`).
					WithArgs(officeID).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if err := tt.sqlor.DeleteOffice(tt.args.officeId); (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.DeleteOffice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLOfficeRepository_GetBuildingAndFloorByOffice(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()
	officeID := uuid.New()

	type args struct {
		officeName string
	}
	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		args      args
		mockSetup func()
		want      uuid.UUID
		want1     int
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "wg",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("wg", 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}).
						AddRow(officeID, "wg", buildingID, 1))
			},
			want:    buildingID,
			want1:   1,
			wantErr: false,
		},
		{
			name: "office_not_found",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "nonexistent",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("nonexistent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    uuid.UUID{},
			want1:   0,
			wantErr: true,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "wg",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("wg", 1).
					WillReturnError(errors.New("database error"))
			},
			want:    uuid.UUID{},
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, got1, err := tt.sqlor.GetBuildingAndFloorByOffice(tt.args.officeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.GetBuildingAndFloorByOffice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLOfficeRepository.GetBuildingAndFloorByOffice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SQLOfficeRepository.GetBuildingAndFloorByOffice() got1 = %v, want %v", got1, tt.want1)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLOfficeRepository_GetOfficesByBuilding(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()
	officeID1 := uuid.New()
	officeID2 := uuid.New()

	expectedOffices := []models.Office{
		{
			OfficeID:    officeID1,
			OfficeName:  "wg",
			BuildingID:  buildingID,
			FloorNumber: 1,
		},
		{
			OfficeID:    officeID2,
			OfficeName:  "office2",
			BuildingID:  buildingID,
			FloorNumber: 2,
		},
	}

	type args struct {
		buildingID string
	}
	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		args      args
		mockSetup func()
		want      []models.Office
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				buildingID: buildingID.String(),
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs(buildingID, "ADMIN_OFFICE").
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}).
						AddRow(officeID1, "wg", buildingID, 1).
						AddRow(officeID2, "office2", buildingID, 2))
			},
			want:    expectedOffices,
			wantErr: false,
		},
		{
			name: "invalid_building_id",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				buildingID: "invalid-uuid",
			},
			mockSetup: func() {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "no_offices_found",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				buildingID: buildingID.String(),
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs(buildingID, "ADMIN_OFFICE").
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}))
			},
			want:    []models.Office{},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				buildingID: buildingID.String(),
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs(buildingID, "ADMIN_OFFICE").
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlor.GetOfficesByBuilding(tt.args.buildingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.GetOfficesByBuilding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLOfficeRepository.GetOfficesByBuilding() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLOfficeRepository_GetAllOffices(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()
	officeID1 := uuid.New()
	officeID2 := uuid.New()

	expectedOffices := []models.Office{
		{
			OfficeID:    officeID1,
			OfficeName:  "wg",
			BuildingID:  buildingID,
			FloorNumber: 1,
		},
		{
			OfficeID:    officeID2,
			OfficeName:  "office2",
			BuildingID:  buildingID,
			FloorNumber: 2,
		},
	}

	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		mockSetup func()
		want      []models.Office
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("ADMIN_OFFICE").
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}).
						AddRow(officeID1, "wg", buildingID, 1).
						AddRow(officeID2, "office2", buildingID, 2))
			},
			want:    expectedOffices,
			wantErr: false,
		},
		{
			name: "no_offices_found",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("ADMIN_OFFICE").
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}))
			},
			want:    []models.Office{},
			wantErr: false,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("ADMIN_OFFICE").
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlor.GetAllOffices()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.GetAllOffices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLOfficeRepository.GetAllOffices() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLOfficeRepository_GetOfficeByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()
	officeID := uuid.New()

	expectedOffice := models.Office{
		OfficeID:    officeID,
		OfficeName:  "wg",
		BuildingID:  buildingID,
		FloorNumber: 1,
	}

	type args struct {
		officeName string
	}
	tests := []struct {
		name      string
		sqlor     *SQLOfficeRepository
		args      args
		mockSetup func()
		want      models.Office
		wantErr   bool
	}{
		{
			name: "success",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "wg",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("wg", 1).
					WillReturnRows(sqlmock.NewRows([]string{"office_id", "office_name", "building_id", "floor_number"}).
						AddRow(officeID, "wg", buildingID, 1))
			},
			want:    expectedOffice,
			wantErr: false,
		},
		{
			name: "office_not_found",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "nonexistent",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("nonexistent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    models.Office{},
			wantErr: true,
		},
		{
			name: "database_error",
			sqlor: &SQLOfficeRepository{
				db: gormDb,
			},
			args: args{
				officeName: "wg",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "offices"`).
					WithArgs("wg", 1).
					WillReturnError(errors.New("database error"))
			},
			want:    models.Office{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := tt.sqlor.GetOfficeByName(tt.args.officeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLOfficeRepository.GetOfficeByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLOfficeRepository.GetOfficeByName() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
