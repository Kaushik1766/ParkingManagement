package floorrepository

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

func TestSQLFloorRepository_AddFloor(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingId  string
		floorNumber int
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
				buildingId:  uuid.New().String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "floors"`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO "slots"`).
					WillReturnResult(sqlmock.NewResult(1, 30))
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
				buildingId:  "invalid-uuid",
				floorNumber: 1,
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
				buildingId:  uuid.New().String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "floors"`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlfr := &SQLFloorRepository{
				db: tt.fields.db,
			}
			if err := sqlfr.AddFloor(tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("SQLFloorRepository.AddFloor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLFloorRepository_DeleteFloor(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingId  string
		floorNumber int
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
				buildingId:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "floor_number"}).AddRow(buildingID, 1))
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "floors"`).
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
				buildingId:  "invalid-uuid",
				floorNumber: 1,
			},
			mockSetup: func() {
			},
			wantErr: true,
		},
		{
			name: "floor not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "database error on find",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "database error on delete",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID.String(),
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "floor_number"}).AddRow(buildingID, 1))
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "floors"`).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlfr := &SQLFloorRepository{
				db: tt.fields.db,
			}
			if err := sqlfr.DeleteFloor(tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("SQLFloorRepository.DeleteFloor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLFloorRepository_GetFloor(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingId  uuid.UUID
		floorNumber int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      int
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "floor_number"}).AddRow(buildingID, 1))
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "floor not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId:  buildingID,
				floorNumber: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID, 1, 1).
					WillReturnError(errors.New("database error"))
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlfr := &SQLFloorRepository{
				db: tt.fields.db,
			}
			got, err := sqlfr.GetFloor(tt.args.buildingId, tt.args.floorNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLFloorRepository.GetFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLFloorRepository.GetFloor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLFloorRepository_GetFloorsByBuildingId(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	buildingID := uuid.New()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		buildingId string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.Floor
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId: buildingID.String(),
			},
			want: []models.Floor{
				{
					BuildingID:  buildingID,
					FloorNumber: 1,
					Slots:       nil,
				},
				{
					BuildingID:  buildingID,
					FloorNumber: 2,
					Slots:       nil,
				},
			},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "floor_number"}).
						AddRow(buildingID, 1).
						AddRow(buildingID, 2))
			},
			wantErr: false,
		},
		{
			name: "invalid building id",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId: "invalid-uuid",
			},
			want: nil,
			mockSetup: func() {
			},
			wantErr: true,
		},
		{
			name: "no floors found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId: buildingID.String(),
			},
			want: []models.Floor{},
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID).
					WillReturnRows(sqlmock.NewRows([]string{"building_id", "floor_number"}))
			},
			wantErr: false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				buildingId: buildingID.String(),
			},
			want: nil,
			mockSetup: func() {
				mock.ExpectQuery(`FROM "floors"`).
					WithArgs(buildingID).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlfr := &SQLFloorRepository{
				db: tt.fields.db,
			}
			got, err := sqlfr.GetFloorsByBuildingId(tt.args.buildingId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLFloorRepository.GetFloorsByBuildingId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLFloorRepository.GetFloorsByBuildingId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSQLFloorRepository(t *testing.T) {
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
		want *SQLFloorRepository
	}{
		{
			name: "success",
			args: args{db: gormDb},
			want: &SQLFloorRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLFloorRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLFloorRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
