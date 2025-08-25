package app_test

import (
	"testing"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/app"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewApp(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://kaushik:123@localhost:5432/kaushik")
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db        *gorm.DB
		want      *app.App
		doesPanic bool
	}{
		{
			name:      "valid gorm.DB",
			db:        func() *gorm.DB { db, _ := db.InitDB(); return db }(),
			want:      &app.App{},
			doesPanic: false,
		},
		{
			name:      "nil db",
			db:        nil,
			want:      &app.App{},
			doesPanic: true,
		},
	}
	for _, tt := range tests {
		if tt.doesPanic {
			assert.Panics(t, func() { app.NewApp(tt.db) })
		} else {
			assert.NotPanics(t, func() { app.NewApp(tt.db) })
		}
	}
}

// func TestApp_Run(t *testing.T) {
// 	t.Setenv("DATABASE_URL", "postgresql://kaushik:123@localhost:5432/kaushik")
// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		db        *gorm.DB
// 		doesPanic bool
// 	}{
// 		{
// 			name:      "valid gorm.DB",
// 			db:        func() *gorm.DB { db, _ := db.InitDB(); return db }(),
// 			doesPanic: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			app := app.NewApp(tt.db)
// 			app.Run()
// 		})
// 	}
// }
