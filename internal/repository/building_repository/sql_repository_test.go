package buildingrepository_test

import (
	"testing"

	"github.com/Kaushik1766/ParkingManagement/db"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	"gorm.io/gorm"
)

func TestSQLBuildingRepository_AddBuilding(t *testing.T) {
	dbURL := ""
	// t.Log("dburl = ", dbURL)
	gormDB, err := db.InitDB(dbURL)
	if err != nil {
		panic("Error connecting to the database: " + err.Error())
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		buildingName string
		wantErr      bool
	}{
		{
			name:         "data input",
			db:           gormDB,
			buildingName: "advant",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlbr := buildingrepository.NewSQLBuildingRepository(tt.db)
			gotErr := sqlbr.AddBuilding(tt.buildingName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddBuilding() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddBuilding() succeeded unexpectedly")
			}
		})
	}
}
