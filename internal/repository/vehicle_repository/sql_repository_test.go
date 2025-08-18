package vehiclerepository_test

import (
	"testing"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db, _ = utils.GetDB()

func TestSQLVehicleRepository_GetVehiclesByUserId(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		userId  uuid.UUID
		want    []models.Vehicle
		wantErr bool
	}{
		{
			name:    "Get vehicle by userid",
			db:      db,
			userId:  uuid.MustParse("9aa170ec-105d-408c-bd8b-b2e66053474c"),
			want:    []models.Vehicle{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlvr := vehiclerepository.NewSQLVehicleRepository(tt.db)
			got, gotErr := sqlvr.GetVehiclesByUserId(tt.userId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetVehiclesByUserId() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetVehiclesByUserId() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetVehiclesByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}
