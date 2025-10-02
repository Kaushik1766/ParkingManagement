package app

import (
	"fmt"
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
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

var (
	userRepo     userrepository.UserStorage                     = nil
	vehicleRepo  vehiclerepository.VehicleStorage               = nil
	officeRepo   officerepository.OfficeStorage                 = nil
	floorRepo    floorrepository.FloorStorage                   = nil
	slotRepo     slotrepository.SlotStorage                     = nil
	buildingRepo buildingrepository.BuildingStorage             = nil
	parkingRepo  parkinghistoryrepository.ParkingHistoryStorage = nil

	userService       userservice.UserManager                 = nil
	assignmentService slotassignment.SlotAssignmentMgr        = nil
	authService       authservice.AuthenticationManager       = nil
	buildingService   buildingservice.BuildingMgr             = nil
	floorService      floorservice.FloorMgr                   = nil
	slotService       slotservice.SlotMgr                     = nil
	officeService     officeservice.OfficeMgr                 = nil
	vehicleService    vehicleservice.VehicleMgr               = nil
	parkingService    parkinghistoryservice.ParkingHistoryMgr = nil
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
	parkingRepo = parkinghistoryrepository.NewSQLParkingRepository(db)

	assignmentService = slotassignment.NewSlotAssignmentService(vehicleRepo, floorRepo, buildingRepo, slotRepo, officeRepo)
	userService = userservice.NewUserService(userRepo, vehicleRepo, officeRepo, assignmentService, buildingRepo)
	authService = authservice.NewAuthService(userRepo)
	buildingService = buildingservice.NewBuildingService(buildingRepo)
	floorService = floorservice.NewFloorService(floorRepo)
	slotService = slotservice.NewSlotService(slotRepo)
	officeService = officeservice.NewOfficeService(officeRepo)
	vehicleService = vehicleservice.NewVehicleService(vehicleRepo, parkingRepo)
	parkingService = parkinghistoryservice.NewParkingHistoryService(parkingRepo, vehicleRepo)

	err := userRepo.(*userrepository.SQLUserRepository).CreateAdminOffice()
	if err != nil {
		color.Red("Error creating admin office: %v", err)
		os.Exit(1)
	}

	// for dev only
	err = userRepo.(*userrepository.SQLUserRepository).SeedBuildingAndOffice()
	if err != nil {
		color.Red("Error seeding test building and offices: %v", err)
		// os.Exit(1)
	}

	err = userRepo.(*userrepository.SQLUserRepository).SeedAdmin()
	if err != nil {
		color.Red("Error seeding admin user: %v", err)
		// os.Exit(1)
	}

	app.AuthHandler = *authhandler.NewWebAuthHandler(authService)
	app.UserHandler = *userhandler.NewWebUserHandler(userService)
	app.BuildingHandler = *buildinghandler.NewWebBuildingHandler(buildingService)
	app.FloorHandler = *floorhandler.NewWebFloorHandler(floorService)
	app.SlotHandler = *slothandler.NewWebSlotHandler(slotService)
	app.OfficeHandler = *officehandler.NewWebOfficeHandler(officeService)
	app.VehicleHandler = *vehiclehandler.NewWebVehicleHandler(vehicleService, userService)
	app.ParkingHandler = *parkinghandler.NewWebParkingHandler(parkingService, vehicleService)

	app.registerRoutes()

	return app
}

func (app *App) Run() {
	fmt.Println("Server started at localhost:3000")

	http.ListenAndServe("localhost:3000", enableCORS(app.apiMux))
}
