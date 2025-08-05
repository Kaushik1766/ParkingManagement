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
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
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
	userService     userservice.UserManager            = nil
	authService     authservice.AuthenticationManager  = nil
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
	userService = userservice.NewUserService(userDb, vehicleDb)
	authService = authservice.NewAuthService(userDb)
	authController = authhandler.NewCliAuthHandler(authService)
	userHandler = userhandler.NewCliUserHandler(userService)

	buildingDb = buildingrepository.NewFileBuildingRepository()
	floorDb = floorrepository.NewFileFloorRepository()
	slotDb = slotrepository.NewFileSlotRepository()

	floorService = floorservice.NewFloorService(floorDb, buildingDb)
	buildingService = buildingservice.NewBuildingService(buildingDb)
	slotService = slotservice.NewSlotService(slotDb, buildingDb, floorDb)
	reader = bufio.NewReader(os.Stdin)
	adminHandler = adminhandler.NewCliAdminHandler(floorService, buildingService, slotService, reader)

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
			color.Yellow("3. Admin Signup")
			color.Yellow("4. Exit")
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
				err := authController.CustomerSignup()
				fmt.Println(err)
			case 3:
				err := authController.AdminSignup()
				fmt.Println(err)
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
				color.Yellow("Enter your choice:")
				color.Yellow("1. Add Building")
				color.Yellow("2. Delete Building")
				color.Yellow("3. Add Floor")
				color.Yellow("4. Delete Floor")
				color.Yellow("5. Add Slots")
				color.Yellow("6. Delete Slots")
				color.Yellow("7. Logout")
				color.Yellow("8. Exit")
				fmt.Scanf("%d", &choice)
				clearScreen()
				switch choice {
				case 1:
					adminHandler.AddBuilding(ctx)
				case 2:
					adminHandler.DeleteBuilding(ctx)
				case 3:
					adminHandler.AddFloors(ctx)
				case 4:
					adminHandler.DeleteFloors(ctx)
				case 5:
					adminHandler.AddSlots(ctx)
				case 6:
					adminHandler.DeleteSlots(ctx)
				case 7:
					ctx = authController.Logout()
				default:
					return
				}

			}
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
