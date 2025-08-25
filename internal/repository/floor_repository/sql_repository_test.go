package floorrepository_test

//
// import (
// 	"testing"
//
// 	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
// 	"github.com/Kaushik1766/ParkingManagement/utils"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )
//
// func TestSQLFloorRepository_AddFloor(t *testing.T) {
// 	db, _ := utils.GetDB()
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		db *gorm.DB
// 		// Named input parameters for target function.
// 		buildingId  uuid.UUID
// 		floorNumber int
// 		wantErr     bool
// 	}{
// 		{
// 			name:        "valid input",
// 			db:          db,
// 			buildingId:  uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
// 			floorNumber: 1,
// 			wantErr:     false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlfr := floorrepository.NewSQLFloorRepository(tt.db)
// 			gotErr := sqlfr.AddFloor(tt.buildingId, tt.floorNumber)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("AddFloor() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("AddFloor() succeeded unexpectedly")
// 			}
// 		})
// 	}
// }
