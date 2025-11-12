package lambda_app

import (
	"os"

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

type LambdaApp struct {
	db *gorm.DB

	userRepo     userrepository.UserStorage
	vehicleRepo  vehiclerepository.VehicleStorage
	officeRepo   officerepository.OfficeStorage
	floorRepo    floorrepository.FloorStorage
	slotRepo     slotrepository.SlotStorage
	buildingRepo buildingrepository.BuildingStorage
	parkingRepo  parkinghistoryrepository.ParkingHistoryStorage

	UserService       userservice.UserManager
	AssignmentService slotassignment.SlotAssignmentMgr
	AuthService       authservice.AuthenticationManager
	BuildingService   buildingservice.BuildingMgr
	FloorService      floorservice.FloorMgr
	SlotService       slotservice.SlotMgr
	OfficeService     officeservice.OfficeMgr
	VehicleService    vehicleservice.VehicleMgr
	ParkingService    parkinghistoryservice.ParkingHistoryMgr
}

func NewLambdaApp(db *gorm.DB) *LambdaApp {
	app := &LambdaApp{
		db: db,
	}

	userRepo := userrepository.NewSQLUserRepository(db)
	vehicleRepo := vehiclerepository.NewSQLVehicleRepository(db)
	officeRepo := officerepository.NewSQLOfficeRepository(db)
	buildingRepo := buildingrepository.NewSQLBuildingRepository(db)
	floorRepo := floorrepository.NewSQLFloorRepository(db)
	slotRepo := slotrepository.NewSQLSlotRepository(db)
	parkingRepo := parkinghistoryrepository.NewSQLParkingRepository(db)

	assignmentService := slotassignment.NewSlotAssignmentService(vehicleRepo, floorRepo, buildingRepo, slotRepo, officeRepo)
	userService := userservice.NewUserService(userRepo, vehicleRepo, officeRepo, assignmentService, buildingRepo)
	authService := authservice.NewAuthService(userRepo)
	buildingService := buildingservice.NewBuildingService(buildingRepo)
	floorService := floorservice.NewFloorService(floorRepo)
	slotService := slotservice.NewSlotService(slotRepo)
	officeService := officeservice.NewOfficeService(officeRepo)
	vehicleService := vehicleservice.NewVehicleService(vehicleRepo, parkingRepo)
	parkingService := parkinghistoryservice.NewParkingHistoryService(parkingRepo, vehicleRepo)

	err := userRepo.CreateAdminOffice()
	if err != nil {
		color.Red("Error creating admin office: %v", err)
		os.Exit(1)
	}

	// for dev only
	err = userRepo.SeedBuildingAndOffice()
	if err != nil {
		color.Red("Error seeding test building and offices: %v", err)
		// os.Exit(1)
	}

	err = userRepo.SeedAdmin()
	if err != nil {
		color.Red("Error seeding admin user: %v", err)
		// os.Exit(1)
	}

	app.UserService = userService
	app.AssignmentService = assignmentService
	app.AuthService = authService
	app.BuildingService = buildingService
	app.FloorService = floorService
	app.SlotService = slotService
	app.OfficeService = officeService
	app.ParkingService = parkingService
	app.VehicleService = vehicleService

	return app
}
