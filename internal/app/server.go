package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/auth_handler"
	buildinghandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/building_handler"
	floorhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/floor_handler"
	officehandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/office_handler"
	parkinghandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/parking_handler"
	slothandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/slot_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/user_handler"
	vehiclehandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/web/vehicle_handler"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

var (
	userRepo     userrepository.UserStorage         = nil
	vehicleRepo  vehiclerepository.VehicleStorage   = nil
	officeRepo   officerepository.OfficeStorage     = nil
	floorRepo    floorrepository.FloorStorage       = nil
	slotRepo     slotrepository.SlotStorage         = nil
	buildingRepo buildingrepository.BuildingStorage = nil

	userService       userservice.UserManager           = nil
	assignmentService slotassignment.SlotAssignmentMgr  = nil
	authService       authservice.AuthenticationManager = nil
)

type App struct {
	db     *gorm.DB
	apiMux *http.ServeMux

	AuthHandler     authhandler.WebAuthHandler
	BuildingHandler buildinghandler.WebBuildingHandler
	FloorHandler    floorhandler.WebFloorHandler
	OfficeHandler   officehandler.WebOfficeHandler
	ParkingHandler  parkinghandler.WebParkingHandler
	SlotHandler     slothandler.WebSlotHandler
	UserHandler     userhandler.WebUserHandler
	VehicleHandler  vehiclehandler.WebVehicleHandler
}

func NewApp(db *gorm.DB) *App {
	app := &App{
		db:     db,
		apiMux: http.NewServeMux(),
	}

	userRepo = userrepository.NewSQLUserRepository(db)
	vehicleRepo = vehiclerepository.NewSQLVehicleRepository(db)
	officeRepo = officerepository.NewSQLOfficeRepository(db)
	buildingRepo = buildingrepository.NewSQLBuildingRepository(db)
	floorRepo = floorrepository.NewSQLFloorRepository(db)
	slotRepo = slotrepository.NewSQLSlotRepository(db)
	officeRepo = officerepository.NewSQLOfficeRepository(db)

	assignmentService = slotassignment.NewSlotAssignmentService(vehicleRepo, floorRepo, buildingRepo, slotRepo, officeRepo)
	userService = userservice.NewUserService(userRepo, vehicleRepo, officeRepo, assignmentService)
	authService = authservice.NewAuthService(userRepo, officeRepo)

	err := userRepo.(*userrepository.SQLUserRepository).CreateAdminOffice()
	if err != nil {
		color.Red("Error creating admin office: %v", err)
		os.Exit(1)
	}
	app.AuthHandler = *authhandler.NewWebAuthHandler(authService)
	app.UserHandler = *userhandler.NewWebUserHandler(userService)

	app.registerRoutes()

	return app
}

func (app *App) Run() {
	fmt.Println("Server started at localhost:3000")
	err := http.ListenAndServe("localhost:3000", app.apiMux)
	if err != nil {
		log.Panic(err)
	}
}
