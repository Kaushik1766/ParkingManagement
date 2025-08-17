package buildingrepository_test

import (
	"slices"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestSQLBuildingRepository_AddBuilding(t *testing.T) {
	db, err := utils.GetDB()
	if err != nil {
		t.Fatal(err)
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
			db:           db,
			buildingName: "advant",
			wantErr:      true,
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

func TestSQLBuildingRepository_GetAllBuildings(t *testing.T) {
	db, err := utils.GetDB()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db      *gorm.DB
		want    []models.Building
		wantErr bool
	}{
		{
			name: "get all buildings",
			db:   db,
			want: []models.Building{
				{
					BuildingID:   uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
					BuildingName: "advant",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlbr := buildingrepository.NewSQLBuildingRepository(tt.db)
			got, gotErr := sqlbr.GetAllBuildings()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetAllBuildings() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetAllBuildings() succeeded unexpectedly")
			}
			for _, b := range tt.want {
				if !slices.ContainsFunc(got, func(g models.Building) bool {
					return g.BuildingID == b.BuildingID && g.BuildingName == b.BuildingName
				}) {
					t.Errorf("GetAllBuildings() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSQLBuildingRepository_GetBuildingByID(t *testing.T) {
	db, _ := utils.GetDB()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		buildingID uuid.UUID
		want       models.Building
		wantErr    bool
	}{
		{
			name:       "get building by ID",
			db:         db,
			buildingID: uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
			want: models.Building{
				BuildingID:   uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
				BuildingName: "advant",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlbr := buildingrepository.NewSQLBuildingRepository(tt.db)
			got, gotErr := sqlbr.GetBuildingByID(tt.buildingID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBuildingByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBuildingByID() succeeded unexpectedly")
			}
			if got.BuildingID != tt.want.BuildingID || got.BuildingName != tt.want.BuildingName {
				t.Errorf("GetBuildingByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLBuildingRepository_GetBuildingByName(t *testing.T) {
	db, _ := utils.GetDB()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db *gorm.DB
		// Named input parameters for target function.
		buildingName string
		want         models.Building
		wantErr      bool
	}{
		{
			name:         "get building by name",
			db:           db,
			buildingName: "advant",
			want: models.Building{
				BuildingID:   uuid.MustParse("26c7d4d8-3ebe-4fad-a3fe-4ef83221acc1"),
				BuildingName: "advant",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlbr := buildingrepository.NewSQLBuildingRepository(tt.db)
			got, gotErr := sqlbr.GetBuildingByName(tt.buildingName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBuildingByName() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBuildingByName() succeeded unexpectedly")
			}
			if got.BuildingID != tt.want.BuildingID || got.BuildingName != tt.want.BuildingName {
				t.Errorf("GetBuildingByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
