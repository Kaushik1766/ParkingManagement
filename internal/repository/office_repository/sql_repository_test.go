package officerepository_test

// import (
// 	"testing"
//
// 	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
// 	"github.com/Kaushik1766/ParkingManagement/utils"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )
//
// var db, _ = utils.GetDB()
//
// func TestSQLOfficeRepository_AddOffice(t *testing.T) {
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		db *gorm.DB
// 		// Named input parameters for target function.
// 		officeName  string
// 		buildingID  uuid.UUID
// 		floorNumber int
// 		wantErr     bool
// 	}{
// 		{
// 			name:        "valid input",
// 			db:          db,
// 			officeName:  "Office A",
// 			buildingID:  uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
// 			floorNumber: 1,
// 			wantErr:     false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sqlor := officerepository.NewSQLOfficeRepository(tt.db)
// 			gotErr := sqlor.AddOffice(tt.officeName, tt.buildingID, tt.floorNumber)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("AddOffice() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("AddOffice() succeeded unexpectedly")
// 			}
// 		})
// 	}
// }
