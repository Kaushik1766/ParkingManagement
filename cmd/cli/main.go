package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
	adminhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/admin_handler"
	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	parkinghandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/parking_handler"
	slotassignmenthandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/slot_assignment_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/user_handler"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	billingservice "github.com/Kaushik1766/ParkingManagement/internal/service/billing_service"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var (
	userDb     userrepository.UserStorage                     = nil
	vehicleDb  vehiclerepository.VehicleStorage               = nil
	floorDb    floorrepository.FloorStorage                   = nil
	slotDb     slotrepository.SlotStorage                     = nil
	buildingDb buildingrepository.BuildingStorage             = nil
	officeDb   officerepository.OfficeStorage                 = nil
	parkingDb  parkinghistoryrepository.ParkingHistoryStorage = nil

	userService           userservice.UserManager                 = nil
	authService           authservice.AuthenticationManager       = nil
	officeService         officeservice.OfficeMgr                 = nil
	floorService          floorservice.FloorMgr                   = nil
	buildingService       buildingservice.BuildingMgr             = nil
	slotService           slotservice.SlotMgr                     = nil
	assignmentService     slotassignment.SlotAssignmentMgr        = nil
	vehicleService        vehicleservice.VehicleMgr               = nil
	parkingHistoryService parkinghistoryservice.ParkingHistoryMgr = nil
	billingService        billingservice.BillingMgr               = nil

	authController        *authhandler.CliAuthHandler                     = nil
	userHandler           *userhandler.CliUserHandler                     = nil
	adminHandler          *adminhandler.CliAdminHandler                   = nil
	slotAssignmentHandler *slotassignmenthandler.CliSlotAssignmentHandler = nil
	parkingHandler        *parkinghandler.CliParkingHandler               = nil

	reader *bufio.Reader   = nil
	ctx    context.Context = nil
)

func init() {
	userDb = userrepository.NewFileUserRepository()
	vehicleDb = vehiclerepository.NewFileVehicleRepository(userDb)
	buildingDb = buildingrepository.NewFileBuildingRepository()
	floorDb = floorrepository.NewFileFloorRepository()
	slotDb = slotrepository.NewFileSlotRepository()
	officeDb = officerepository.NewFileOfficeRepository()
	parkingDb = parkinghistoryrepository.NewFileParkingHistoryRepository()

	authService = authservice.NewAuthService(userDb, officeDb)
	floorService = floorservice.NewFloorService(floorDb, buildingDb)
	buildingService = buildingservice.NewBuildingService(buildingDb)
	slotService = slotservice.NewSlotService(slotDb, buildingDb, floorDb)
	officeService = officeservice.NewOfficeService(officeDb, buildingDb, floorDb)
	assignmentService = slotassignment.NewSlotAssignmentService(vehicleDb, floorDb, buildingDb, slotDb, officeDb)
	userService = userservice.NewUserService(userDb, vehicleDb, officeDb, assignmentService)
	vehicleService = vehicleservice.NewVehicleService(vehicleDb, parkingDb)
	parkingHistoryService = parkinghistoryservice.NewParkingHistoryService(parkingDb, vehicleDb)
	billingService = billingservice.NewBillingService(userService, parkingHistoryService)

	reader = bufio.NewReader(os.Stdin)

	authController = authhandler.NewCliAuthHandler(authService, officeService)
	userHandler = userhandler.NewCliUserHandler(userService)
	adminHandler = adminhandler.NewCliAdminHandler(floorService, buildingService, slotService, reader, officeService)
	slotAssignmentHandler = slotassignmenthandler.NewCliSlotAssignmentHandler(assignmentService, userService, slotService, officeService)
	parkingHandler = parkinghandler.NewCliParkingHandler(vehicleService, userService, parkingHistoryService)

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
	officeDb.(*officerepository.FileOfficeRepository).SerializeData()
	parkingDb.(*parkinghistoryrepository.FileParkingRepository).SerializeData()
}

func main() {
	defer cleanup()
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(logFile)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)

	go func() {
		<-ch
		cleanup()
		os.Exit(0)
	}()

	go func() {
		for {
			billingService.GenerateMonthlyInvoice()
		}
	}()
	var choice int
	for {
		if ctx == nil {
			clearScreen()
			// color.Cyan(menuconstants.WelcomeMessage)
			(figure.NewColorFigure("1337Park", "", "blue", true)).Print()
			fmt.Println("")
			color.Yellow(menuconstants.LoginOption)
			color.Yellow(menuconstants.SignupOption)
			color.Yellow(menuconstants.AdminSignupOption)
			color.Yellow(menuconstants.ExitOption)
			fmt.Scanf("%d", &choice)
			clearScreen()
			switch choice {
			case 1:
				ctx = context.Background()
				userCtx, err := authController.Login(ctx)
				if err != nil {
					ctx = nil
					log.Println("Login failed:", err)
				} else {
					ctx = userCtx
				}
			case 2:
				authController.CustomerSignup()
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
				(figure.NewColorFigure("Customer Menu", "", "blue", true)).Print()
				fmt.Println("")
				color.Cyan(menuconstants.EnterYourChoice)
				color.Yellow(menuconstants.CustomerProfileMenu)
				color.Yellow(menuconstants.CustomerRegistrationMenu)
				color.Yellow(menuconstants.CustomerParkingMenu)
				color.Yellow(menuconstants.CustomerLogout)
				color.Yellow(menuconstants.CustomerExit)
				fmt.Scanf("%d", &choice)
				clearScreen()
				switch choice {
				case 1:
					profileManagement()
				case 2:
					registrationMenu()
				case 3:
					parkingMenu()
				case 4:
					ctx = authController.Logout()
				default:
					return
				}
			} else {
				clearScreen()
				// color.Cyan(menuconstants.AdminPageTitle)
				(figure.NewColorFigure("Admin Menu", "", "blue", true)).Print()
				fmt.Println("")
				color.Cyan(menuconstants.AdminBuildingManagement)
				color.Cyan(menuconstants.AdminFloorManagement)
				color.Cyan(menuconstants.AdminSlotManagement)
				color.Cyan(menuconstants.AdminOfficeManagement)
				color.Cyan(menuconstants.AdminUnassignedSlotMgmt)
				color.Cyan(menuconstants.AdminLogout)
				color.Cyan(menuconstants.AdminExit)

				color.Yellow(menuconstants.EnterYourChoice)
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
					unassignedSlotManagement()
				case 6:
					ctx = authController.Logout()
				default:
					return
				}

			}
		}
	}
}

func parkingMenu() {
	for {
		clearScreen()
		(figure.NewColorFigure("Parking Menu", "", "blue", true)).Print()
		fmt.Println("")
		color.Cyan(menuconstants.ParkingMenuTitle)
		color.Yellow(menuconstants.ParkVehicle)
		color.Yellow(menuconstants.UnparkVehicle)
		color.Yellow(menuconstants.ViewParkings)
		color.Yellow(menuconstants.ParkingExit)
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			parkingHandler.Park(ctx)
		case 2:
			parkingHandler.Unpark(ctx)
		case 3:
			parkingHandler.ViewParkingHistory(ctx)
		default:
			return
		}
	}
}

func unassignedSlotManagement() {
	for {
		clearScreen()
		(figure.NewColorFigure("Assignment menu", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.ViewUnassignedSlots)
		color.Yellow(menuconstants.AssignSlot)
		color.Yellow(menuconstants.UnassignedSlotExit)
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			slotAssignmentHandler.ViewVehiclesWithUnassignedSlots(ctx)
		case 2:
			slotAssignmentHandler.AssignSlot(ctx)
		default:
			return
		}
	}
}

func buildingManagement() {
	for {
		clearScreen()
		(figure.NewColorFigure("Building Management", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.AddBuilding)
		color.Yellow(menuconstants.DeleteBuilding)
		color.Yellow(menuconstants.ListBuildings)
		color.Yellow(menuconstants.BuildingExit)
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
		clearScreen()
		(figure.NewColorFigure("Floor Management", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.AddFloor)
		color.Yellow(menuconstants.DeleteFloor)
		color.Yellow(menuconstants.ListFloors)
		color.Yellow(menuconstants.FloorExit)
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
		clearScreen()
		(figure.NewColorFigure("Slot Management", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.AddSlots)
		color.Yellow(menuconstants.DeleteSlots)
		color.Yellow(menuconstants.ViewSlots)
		color.Yellow(menuconstants.SlotExit)
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
		clearScreen()
		(figure.NewColorFigure("Office Management", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.AddOffice)
		color.Yellow(menuconstants.RemoveOffice)
		color.Yellow(menuconstants.ListOffices)
		color.Yellow(menuconstants.OfficeExit)
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

func profileManagement() {
	for {
		clearScreen()
		(figure.NewColorFigure("Profile Management", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.ProfileUpdate)
		color.Yellow(menuconstants.ProfileView)
		color.Yellow(menuconstants.ProfileExit)
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			userHandler.UpdateProfile(ctx)
		case 2:
			userHandler.GetUserProfile(ctx)
		default:
			return
		}
	}
}

func registrationMenu() {
	for {
		clearScreen()
		(figure.NewColorFigure("Registration Menu", "", "blue", true)).Print()
		fmt.Println("")
		color.Yellow(menuconstants.EnterYourChoice)
		color.Yellow(menuconstants.RegistrationAddVehicle)
		color.Yellow(menuconstants.RegistrationViewVehicles)
		color.Yellow(menuconstants.RegistrationUnregisterVehicle)
		color.Yellow(menuconstants.RegistrationExit)
		var choice int
		fmt.Scanf("%d", &choice)
		clearScreen()
		switch choice {
		case 1:
			userHandler.RegisterVehicle(ctx)
		case 2:
			userHandler.GetRegisteredVehicles(ctx)
		case 3:
			userHandler.UnregisterVehicle(ctx)
		default:
			return
		}
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
