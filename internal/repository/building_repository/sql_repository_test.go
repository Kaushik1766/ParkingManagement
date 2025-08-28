package buildingrepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewSQLBuildingRepository(t *testing.T) {
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
		want *SQLBuildingRepository
	}{
		{
			name: "success",
			args: args{db: gormDb},
			want: &SQLBuildingRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLBuildingRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLBuildingRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLBuildingRepository_AddBuilding(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: "advant",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "buildings"`).
					WithArgs("advant").
					WillReturnRows(sqlmock.NewRows([]string{"building_id"}).AddRow(uuid.New().String()))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "duplicate building",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: "advant",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "buildings"`).
					WithArgs("advant").
					WillReturnError(errors.New("duplicate building"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "admin building",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: constants.AdminBuilding,
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlbr := &SQLBuildingRepository{
				db: tt.fields.db,
			}
			if err := sqlbr.AddBuilding(tt.args.buildingName); (err != nil) != tt.wantErr {
				t.Errorf("AddBuilding() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLBuildingRepository_DeleteBuildingByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingID string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: uuid.New().String(),
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "buildings"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "invalid building id",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: "asdfasdf",
			},
			mockSetup: func() {

			},
			wantErr: true,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: uuid.New().String(),
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "buildings"`).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlbr := &SQLBuildingRepository{
				db: tt.fields.db,
			}
			if err := sqlbr.DeleteBuildingByID(tt.args.buildingID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBuildingByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLBuildingRepository_GetAllBuildings(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		want      []models.Building
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			want: []models.Building{
				{
					BuildingID:   uuid.Nil,
					BuildingName: "advant",
					Floors:       nil,
				},
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs(constants.AdminBuilding).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "building_name"}).
						AddRow(uuid.Nil, "advant"))
			},
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			want: []models.Building{},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs(constants.AdminBuilding).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlbr := &SQLBuildingRepository{
				db: tt.fields.db,
			}
			got, err := sqlbr.GetAllBuildings()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBuildings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllBuildings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLBuildingRepository_GetBuildingByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	wantBuilding := models.Building{
		BuildingID:   uuid.New(),
		BuildingName: "advant",
		Floors:       nil,
	}

	notFoundID := uuid.New()
	dbErrorID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingID uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.Building
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: wantBuilding.BuildingID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs(wantBuilding.BuildingID, constants.AdminBuilding, 1).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "building_name"}).AddRow(wantBuilding.BuildingID, wantBuilding.BuildingName))
			},
			wantErr: false,
			want:    wantBuilding,
		},
		{
			name: "building not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: notFoundID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs(notFoundID, constants.AdminBuilding, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			want:    models.Building{},
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingID: dbErrorID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs(dbErrorID, constants.AdminBuilding, 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
			want:    models.Building{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlbr := &SQLBuildingRepository{
				db: tt.fields.db,
			}
			got, err := sqlbr.GetBuildingByID(tt.args.buildingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBuildingByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBuildingByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLBuildingRepository_GetBuildingByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))

	wantBuilding := models.Building{
		BuildingID:   uuid.New(),
		BuildingName: "advant",
		Floors:       nil,
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.Building
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: "advant",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs("advant", 1).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "building_name"}).AddRow(wantBuilding.BuildingID, wantBuilding.BuildingName))
			},
			wantErr: false,
			want:    wantBuilding,
		},
		{
			name: "admin building",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: constants.AdminBuilding,
			},
			mockSetup: func() {
			},
			wantErr: true,
			want:    models.Building{},
		},
		{
			name: "building not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: "nonexistent",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs("nonexistent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			want:    models.Building{},
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingName: "advant",
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "buildings"`).
					WithArgs("advant", 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
			want:    models.Building{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlbr := &SQLBuildingRepository{
				db: tt.fields.db,
			}
			got, err := sqlbr.GetBuildingByName(tt.args.buildingName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBuildingByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBuildingByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
