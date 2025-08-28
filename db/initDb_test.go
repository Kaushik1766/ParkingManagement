package db

import (
	"os"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"gorm.io/gorm"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "missing database url",
			wantErr: true,
		},
		{
			name:    "invalid database url",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var
			originalURL := os.Getenv("DATABASE_URL")
			defer func() {
				if originalURL != "" {
					_ = os.Setenv("DATABASE_URL", originalURL)
				} else {
					_ = os.Unsetenv("DATABASE_URL")
				}
			}()

			// Set up test environment
			if tt.name == "missing database url" {
				_ = os.Unsetenv("DATABASE_URL")
			} else if tt.name == "invalid database url" {
				_ = os.Setenv("DATABASE_URL", "invalid-url")
			}

			got, err := InitDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Only check that we get a database object when there's no error
			if !tt.wantErr && got == nil {
				t.Errorf("InitDB() returned nil database when expecting success")
			}
		})
	}
}

func TestMigrateModels(t *testing.T) {
	type args struct {
		db     *gorm.DB
		models []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil database",
			args: args{
				db:     nil,
				models: []any{&models.User{}},
			},
			wantErr: false,
		},
		{
			name: "empty models",
			args: args{
				db:     nil,
				models: []any{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MigrateModels(tt.args.db, tt.args.models...); (err != nil) != tt.wantErr {
				t.Errorf("MigrateModels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
