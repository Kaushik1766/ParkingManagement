package db_test

import (
	"testing"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"gorm.io/gorm"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    *gorm.DB
		wantErr bool
		env     string
	}{
		{
			name:    "valid db url",
			want:    &gorm.DB{},
			wantErr: false,
			env:     "postgresql://kaushik:123@localhost:5432/kaushik",
		},
		{
			name:    "invalid db url",
			want:    nil,
			wantErr: true,
			env:     "invalid_url",
		},
	}
	for _, tt := range tests {
		t.Setenv("DATABASE_URL", tt.env)
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := db.InitDB()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InitDB() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InitDB() succeeded unexpectedly")
			}
			if got.Error != nil {
				t.Errorf("InitDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigrateModels(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://kaushik:123@localhost:5432/kaushik")
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db      *gorm.DB
		models  []any
		wantErr bool
	}{
		{
			name:    "nil db",
			db:      nil,
			models:  []any{},
			wantErr: false,
		},
		{
			name: "valid db and models",
			db:   func() *gorm.DB { db, _ := db.InitDB(); return db }(),
			models: []any{
				models.Building{},
				models.Floor{},
				models.Slot{},
				models.Office{},
				models.User{},
				models.Vehicle{},
				models.ParkingHistory{},
			},
			wantErr: false,
		},
		{
			name: "invalid relations",
			db:   func() *gorm.DB { db, _ := db.InitDB(); return db }(),
			models: []any{
				struct{ ID int }{},
				struct {
					xyz int `gorm:"foreignKey:NonExistentID"`
				}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := db.MigrateModels(tt.db, tt.models...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MigrateModels() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MigrateModels() succeeded unexpectedly")
			}
		})
	}
}
