package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	adminhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/admin_handler"
	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/user_handler"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	"github.com/fatih/color"
)

var (
	userDb          userrepository.UserStorage         = nil
	vehicleDb       vehiclerepository.VehicleStorage   = nil
	floorDb         floorrepository.FloorStorage       = nil
	slotDb          slotrepository.SlotStorage         = nil
	buildingDb      buildingrepository.BuildingStorage = nil
	officeDb        officerepository.OfficeStorage     = nil
	userService     userservice.UserManager            = nil
	authService     authservice.AuthenticationManager  = nil
	officeService   officeservice.OfficeMgr            = nil
	authController  *authhandler.CliAuthHandler        = nil
	userHandler     *userhandler.CliUserHandler        = nil
	adminHandler    *adminhandler.CliAdminHandler      = nil
	floorService    floorservice.FloorMgr              = nil
	buildingService buildingservice.BuildingMgr        = nil
	slotService     slotservice.SlotMgr                = nil
	reader          *bufio.Reader                      = nil
	ctx             context.Context                    = nil
)

func init() {
	userDb = userrepository.NewFileUserRepository()
	vehicleDb = vehiclerepository.NewFileVehicleRepository(userDb)
	buildingDb = buildingrepository.NewFileBuildingRepository()
	floorDb = floorrepository.NewFileFloorRepository()
	slotDb = slotrepository.NewFileSlotRepository()
	officeDb = officerepository.NewFileOfficeRepository()

	authService = authservice.NewAuthService(userDb, officeDb)
	floorService = floorservice.NewFloorService(floorDb, buildingDb)
	buildingService = buildingservice.NewBuildingService(buildingDb)
	slotService = slotservice.NewSlotService(slotDb, buildingDb, floorDb)
	officeService = officeservice.NewOfficeService(officeDb, buildingDb, floorDb)
	userService = userservice.NewUserService(userDb, vehicleDb, officeDb)

	authController = authhandler.NewCliAuthHandler(authService, officeService)
	userHandler = userhandler.NewCliUserHandler(userService)
	adminHandler = adminhandler.NewCliAdminHandler(floorService, buildingService, slotService, reader, officeService)

	reader = bufio.NewReader(os.Stdin)

	loadLogin()
}

func loadLogin() {
	data, _ := os.ReadFile("token.txt")
	token := string(data)
	ctx, _ = authenticationmiddleware.CliAuthenticate(context.Background(), token)
}

func cleanup() {
	color.Green("Cleaning up...")
	userDb.(*userrepository.FileUserRepository).SerializeData()
	vehicleDb.(*vehiclerepository.FileVehicleRepository).SerializeData()
	floorDb.(*floorrepository.FileFloorRepository).SerializeData()
	slotDb.(*slotrepository.FileSlotRepository).SerializeData()
	buildingDb.(*buildingrepository.FileBuildingRepository).SerializeData()
}

func main() {
	defer cleanup()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)

	go func() {
		<-ch
		cleanup()
		os.Exit(0)
	}()
	var choice int
	for {
		if ctx == nil {
			clearScreen()
			color.Cyan("Welcome to Parking Management System")
			color.Yellow("1. Login")
			color.Yellow("2. Signup")
			color.Yellow("3. Exit")
			fmt.Scanf("%d", &choice)
			clearScreen()
			switch choice {
			case 1:
				ctx = context.Background()
				userCtx, err := authController.Login(ctx)
				if err != nil {
					color.Red("Error during login: %v", err)
					color.Cyan("Please try again or signup if you don't have an account.")
					ctx = nil
					fmt.Scanln()
				} else {
					ctx = userCtx
				}
			case 2:
				authController.CustomerSignup()
			default:
				return
			}
		} else {
			user, ok := ctx.Value(constants.User).(userjwt.UserJwt)
			if !ok {
				continue
			}
			if user.Role == roles.Customer {
				clearScreen()
				color.Cyan("Enter your choice:")
				color.Yellow("1. Update profile")
				color.Yellow("2. Register vehicle")
				color.Yellow("3. View Profile")
				color.Yellow("4. View Registered Vehicles")
				color.Yellow("5. Logout")
				color.Yellow("6. Exit")
				fmt.Scanf("%d", &choice)
				clearScreen()
				switch choice {
				case 1:
					userHandler.UpdateProfile(ctx)
				case 2:
					userHandler.RegisterVehicle(ctx)
				case 3:
					userHandler.GetUserProfile(ctx)
				case 4:
					userHandler.GetRegisteredVehicles(ctx)
				case 5:
					ctx = authController.Logout()
				default:
					return
				}
			} else {
				clearScreen()
				color.Cyan("Admin page: ")
				color.Cyan("1. Building Management")
				color.Cyan("2. Floor Management")
				color.Cyan("3. Slot Management")
				color.Cyan("4. Office Management")
				color.Cyan("5. Logout")
				color.Cyan("6. Exit")

				color.Yellow("Enter your choice:")
				fmt.Scanf("%d", &choice)
				clearScreen()
				switch choice {
				case 1:
					buildingManagement()
				case 2:
					floorManagement()
				case 3:
					slotManagement()
				case 4:
					officeManagement()
				case 5:
					ctx = authController.Logout()
				default:
					return
				}

			}
		}
	}
}

func buildingManagement() {
	for {
		color.Yellow("Enter your choice:")
		color.Yellow("1. Add Building")
		color.Yellow("2. Delete Building")
		color.Yellow("3. List Buildings")
		color.Yellow("4. Exit")
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			adminHandler.AddBuilding(ctx)
		case 2:
			adminHandler.DeleteBuilding(ctx)
		case 3:
			adminHandler.ListBuildings(ctx)
		default:
			return
		}
	}
}

func floorManagement() {
	for {
		color.Yellow("Enter your choice:")
		color.Yellow("1. Add Floor")
		color.Yellow("2. Delete Floor")
		color.Yellow("3. List Floors")
		color.Yellow("4. Exit")
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			adminHandler.AddFloors(ctx)
		case 2:
			adminHandler.DeleteFloors(ctx)
		case 3:
			adminHandler.ListFloors(ctx)
		default:
			return
		}
	}
}

func slotManagement() {
	for {
		color.Yellow("Enter your choice:")
		color.Yellow("1. Add Slots")
		color.Yellow("2. Delete Slots")
		color.Yellow("3. View Slots")
		color.Yellow("4. Exit")
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			adminHandler.AddSlots(ctx)
		case 2:
			adminHandler.DeleteSlots(ctx)
		case 3:
			adminHandler.ListSlots(ctx)
		default:
			return
		}
	}
}

func officeManagement() {
	for {
		color.Yellow("Enter your choice:")
		color.Yellow("1. Add Office")
		color.Yellow("2. Remove Office")
		color.Yellow("3. List Offices")
		color.Yellow("4. Exit")
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			adminHandler.AddOffice(ctx)
		case 2:
			adminHandler.RemoveOffice(ctx)
		case 3:
			adminHandler.ListOffices(ctx)
		default:
			return
		}
	}
}

func clearScreen() {
	// var cmd *exec.Cmd
	// switch runtime.GOOS {
	// case "windows":
	// 	cmd = exec.Command("cmd", "/c", "cls")
	// default:
	// 	cmd = exec.Command("clear")
	// }
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}
