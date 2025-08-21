package app

import (
	"net/http"

	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/auth_handler"
	buildinghandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/building_handler"
	floorhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/floor_handler"
	officehandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/office_handler"
	parkinghandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/parking_handler"
	slothandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/slot_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/user_handler"
	vehiclehandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/vehicle_handler"
	"gorm.io/gorm"
)

type App struct {
	db  *gorm.DB
	mux *http.ServeMux

	AuthHandler     authhandler.WebAuthHandler
	BuildingHandler buildinghandler.WebBuildingHandler
	FloorHandler    floorhandler.WebFloorHandler
	OfficeHandler   officehandler.WebOfficeHandler
	ParkingHandler  parkinghandler.WebParkingHandler
	SlotHandler     slothandler.WebSlotHandler
	UserHandler     userhandler.WebUserHandler
	VehicleHandler  vehiclehandler.WebVehicleHandler
}

func NewApp(db *gorm.DB, mux *http.ServeMux) *App {
	return &App{
		db:  db,
		mux: mux,

		// AuthHandler:     authhandler.NewWebAuthHandler(db),
		// BuildingHandler: buildinghandler.NewWebBuildingHandler(db),
		// FloorHandler:    floorhandler.NewWebFloorHandler(db),
		// OfficeHandler:   officehandler.NewWebOfficeHandler(db),
		// ParkingHandler:  parkinghandler.NewWebParkingHandler(db),
		// SlotHandler:     slothandler.NewWebSlotHandler(db),
		// UserHandler:     userhandler.NewWebUserHandler(db),
		// VehicleHandler:  vehiclehandler.NewWebVehicleHandler(db),
	}
}
