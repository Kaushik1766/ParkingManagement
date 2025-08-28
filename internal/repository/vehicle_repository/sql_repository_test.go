package vehiclerepository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	roles "github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSQLVehicleRepository_GetVehiclesByUserId(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	userID := uuid.New()
	officeID := uuid.New()
	expectedVehicles := []models.Vehicle{
		{
			VehicleID:   uuid.New(),
			NumberPlate: "ABC123",
			UserID:      userID,
			User: models.User{
				UserID:   userID,
				Name:     "kaushik",
				Email:    "kaushik@a.com",
				Password: "123",
				Role:     roles.Customer,
				IsActive: true,
				OfficeID: officeID,
			},
			VehicleType: vehicletypes.TwoWheeler,
			IsActive:    true,
		},
		{
			VehicleID:   uuid.New(),
			NumberPlate: "XYZ789",
			UserID:      userID,
			User: models.User{
				UserID:   userID,
				Name:     "kaushik",
				Email:    "kaushik@a.com",
				Password: "123",
				Role:     roles.Customer,
				IsActive: true,
				OfficeID: officeID,
			},
			VehicleType: vehicletypes.FourWheeler,
			IsActive:    true,
		},
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		userId uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      []models.Vehicle
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				userId: userID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE user_id = \$1`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "number_plate", "vehicle_type", "user_id", "is_active", "assigned_slot_id"}).
						AddRow(expectedVehicles[0].VehicleID, "ABC123", vehicletypes.TwoWheeler, userID, true, nil).
						AddRow(expectedVehicles[1].VehicleID, "XYZ789", vehicletypes.FourWheeler, userID, true, nil))
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."user_id" = \$1`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "email", "password", "role", "is_active", "office_id"}).
						AddRow(userID, "kaushik", "kaushik@a.com", "123", roles.Customer, true, officeID))
			},
			want:    expectedVehicles,
			wantErr: false,
		},
		{
			name: "database_error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				userId: userID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE user_id = \$1`).
					WithArgs(userID).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			got, err := sqlvr.GetVehiclesByUserId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.GetVehiclesByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLVehicleRepository.GetVehiclesByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLVehicleRepository_RemoveVehicle(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		numberplate string
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
				numberplate: "ABC123",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "vehicles" WHERE number_plate = \$1`).
					WithArgs("ABC123").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "ABC123",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "vehicles" WHERE number_plate = \$1`).
					WithArgs("ABC123").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			if err := sqlvr.RemoveVehicle(tt.args.numberplate); (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.RemoveVehicle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLVehicleRepository_GetVehicleById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	vehicleID := uuid.New()
	expectedVehicle := models.Vehicle{
		VehicleID:   vehicleID,
		NumberPlate: "ABC123",
		UserID:      uuid.New(),
		VehicleType: vehicletypes.TwoWheeler,
		IsActive:    true,
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		vehicleId uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.Vehicle
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				vehicleId: vehicleID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE vehicle_id = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs(vehicleID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "number_plate", "vehicle_type", "user_id", "is_active"}).
						AddRow(vehicleID, "ABC123", vehicletypes.TwoWheeler, expectedVehicle.UserID, true))
			},
			want:    expectedVehicle,
			wantErr: false,
		},
		{
			name: "vehicle not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				vehicleId: vehicleID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE vehicle_id = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs(vehicleID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    models.Vehicle{},
			wantErr: true,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				vehicleId: vehicleID,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE vehicle_id = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs(vehicleID, 1).
					WillReturnError(errors.New("database error"))
			},
			want:    models.Vehicle{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			got, err := sqlvr.GetVehicleById(tt.args.vehicleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.GetVehicleById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.NumberPlate, tt.want.NumberPlate) {
				t.Errorf("SQLVehicleRepository.GetVehicleById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLVehicleRepository_GetVehicleByNumberPlate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	vehicleID := uuid.New()
	userID := uuid.New()
	expectedVehicle := models.Vehicle{
		VehicleID:   vehicleID,
		NumberPlate: "ABC123",
		UserID:      userID,
		User: models.User{
			UserID:   userID,
			Name:     "kaushik",
			Email:    "kaushik@a.com",
			Password: "123",
			Role:     roles.Customer,
			IsActive: true,
		},
		VehicleType: vehicletypes.TwoWheeler,
		IsActive:    true,
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		numberplate string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.Vehicle
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "ABC123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE number_plate = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs("ABC123", 1).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "number_plate", "vehicle_type", "user_id", "is_active", "assigned_slot_id"}).
						AddRow(vehicleID, "ABC123", vehicletypes.TwoWheeler, userID, true, nil))
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."user_id" = \$1`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "email", "password", "role", "is_active", "office_id"}).
						AddRow(userID, "kaushik", "kaushik@a.com", "123", roles.Customer, true, uuid.Nil))
			},
			want:    expectedVehicle,
			wantErr: false,
		},
		{
			name: "vehicle not found",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "NONEXISTENT",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE number_plate = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs("NONEXISTENT", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    models.Vehicle{},
			wantErr: true,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "ABC123",
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE number_plate = \$1 ORDER BY "vehicles"\."vehicle_id" LIMIT \$2`).
					WithArgs("ABC123", 1).
					WillReturnError(errors.New("database error"))
			},
			want:    models.Vehicle{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			got, err := sqlvr.GetVehicleByNumberPlate(tt.args.numberplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.GetVehicleByNumberPlate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.NumberPlate, tt.want.NumberPlate) {
				t.Errorf("SQLVehicleRepository.GetVehicleByNumberPlate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLVehicleRepository_GetVehiclesWithUnassignedSlots(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	expectedVehicles := []models.Vehicle{
		{
			VehicleID:   uuid.New(),
			NumberPlate: "ABC123",
			UserID:      uuid.New(),
			VehicleType: vehicletypes.TwoWheeler,
			IsActive:    true,
		},
		{
			VehicleID:   uuid.New(),
			NumberPlate: "XYZ789",
			UserID:      uuid.New(),
			VehicleType: vehicletypes.FourWheeler,
			IsActive:    true,
		},
	}

	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name         string
		fields       fields
		wantVehicles []models.Vehicle
		mockSetup    func()
		wantErr      bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE assigned_building_id IS NULL AND assigned_floor_number IS NULL AND assigned_slot_number IS NULL`).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "number_plate", "vehicle_type", "user_id", "is_active"}).
						AddRow(expectedVehicles[0].VehicleID, "ABC123", vehicletypes.TwoWheeler, expectedVehicles[0].UserID, true).
						AddRow(expectedVehicles[1].VehicleID, "XYZ789", vehicletypes.FourWheeler, expectedVehicles[1].UserID, true))
			},
			wantVehicles: expectedVehicles,
			wantErr:      false,
		},
		{
			name: "no vehicles found",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE assigned_building_id IS NULL AND assigned_floor_number IS NULL AND assigned_slot_number IS NULL`).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "number_plate", "vehicle_type", "user_id", "is_active"}))
			},
			wantVehicles: []models.Vehicle{},
			wantErr:      false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "vehicles" WHERE assigned_building_id IS NULL AND assigned_floor_number IS NULL AND assigned_slot_number IS NULL`).
					WillReturnError(errors.New("database error"))
			},
			wantVehicles: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			gotVehicles, err := sqlvr.GetVehiclesWithUnassignedSlots()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.GetVehiclesWithUnassignedSlots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(gotVehicles) != len(tt.wantVehicles) {
				t.Errorf("SQLVehicleRepository.GetVehiclesWithUnassignedSlots() = %v, want %v", gotVehicles, tt.wantVehicles)
			}
		})
	}
}

func TestSQLVehicleRepository_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	vehicle := models.Vehicle{
		VehicleID:   uuid.New(),
		NumberPlate: "ABC123",
		UserID:      uuid.New(),
		VehicleType: vehicletypes.TwoWheeler,
		IsActive:    true,
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		vehicle models.Vehicle
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
				vehicle: vehicle,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "vehicles" SET "number_plate"=\$1,"vehicle_type"=\$2,"user_id"=\$3,"assigned_building_id"=\$4,"assigned_floor_number"=\$5,"assigned_slot_number"=\$6,"is_active"=\$7 WHERE "vehicle_id" = \$8`).
					WithArgs("ABC123", vehicletypes.TwoWheeler, vehicle.UserID, nil, nil, nil, true, vehicle.VehicleID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				vehicle: vehicle,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "vehicles" SET "number_plate"=\$1,"vehicle_type"=\$2,"user_id"=\$3,"assigned_building_id"=\$4,"assigned_floor_number"=\$5,"assigned_slot_number"=\$6,"is_active"=\$7 WHERE "vehicle_id" = \$8`).
					WithArgs("ABC123", vehicletypes.TwoWheeler, vehicle.UserID, nil, nil, nil, true, vehicle.VehicleID).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			if err := sqlvr.Save(tt.args.vehicle); (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSQLVehicleRepository(t *testing.T) {
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
		want *SQLVehicleRepository
	}{
		{
			name: "success",
			args: args{db: gormDb},
			want: &SQLVehicleRepository{
				db: gormDb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSQLVehicleRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLVehicleRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLVehicleRepository_AddVehicle(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	userID := uuid.New()
	expectedVehicle := models.Vehicle{
		VehicleID:   uuid.New(),
		NumberPlate: "ABC123",
		UserID:      userID,
		VehicleType: vehicletypes.TwoWheeler,
		IsActive:    true,
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		numberplate string
		userid      uuid.UUID
		vehicleType vehicletypes.VehicleType
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      models.Vehicle
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "ABC123",
				userid:      userID,
				vehicleType: vehicletypes.TwoWheeler,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "vehicles" \("number_plate","vehicle_type","user_id","is_active"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "vehicle_id","assigned_building_id","assigned_floor_number","assigned_slot_number"`).
					WithArgs("ABC123", vehicletypes.TwoWheeler, userID, true).
					WillReturnRows(sqlmock.NewRows([]string{"vehicle_id", "assigned_building_id", "assigned_floor_number", "assigned_slot_number"}).
						AddRow(expectedVehicle.VehicleID, nil, nil, nil))
				mock.ExpectCommit()
			},
			want:    expectedVehicle,
			wantErr: false,
		},
		{
			name: "database_error",
			fields: fields{
				db: gormDb,
			},
			args: args{
				numberplate: "ABC123",
				userid:      userID,
				vehicleType: vehicletypes.TwoWheeler,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "vehicles" \("number_plate","vehicle_type","user_id","is_active"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "vehicle_id","assigned_building_id","assigned_floor_number","assigned_slot_number"`).
					WithArgs("ABC123", vehicletypes.TwoWheeler, userID, true).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			want:    models.Vehicle{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sqlvr := &SQLVehicleRepository{
				db: tt.fields.db,
			}
			got, err := sqlvr.AddVehicle(tt.args.numberplate, tt.args.userid, tt.args.vehicleType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLVehicleRepository.AddVehicle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.NumberPlate != tt.want.NumberPlate {
				t.Errorf("SQLVehicleRepository.AddVehicle() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
