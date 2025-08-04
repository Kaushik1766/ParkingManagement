package adminhandler

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	"github.com/fatih/color"
)

type CliAdminHandler struct {
	floorService    floorservice.FloorMgr
	buildingService buildingservice.BuildingMgr
	slotService     slotservice.SlotMgr
}

func NewCliAdminHandler(
	floorService floorservice.FloorMgr,
	buildingService buildingservice.BuildingMgr,
	slotService slotservice.SlotMgr,
) *CliAdminHandler {
	return &CliAdminHandler{
		floorService:    floorService,
		buildingService: buildingService,
		slotService:     slotService,
	}
}

func (h *CliAdminHandler) AddBuilding(ctx context.Context) {
	color.Blue("Enter the name of the building to add:")
	var buildingName string
	fmt.Scanln(&buildingName)
	err := h.buildingService.AddBuilding(ctx, buildingName)
	if err != nil {
		color.Red("Error adding building: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
	}
}

func (h *CliAdminHandler) DeleteBuilding(ctx context.Context) {
	color.Blue("Enter the name of the building to delete:")
	var buildingName string
	fmt.Scanln(&buildingName)
	err := h.buildingService.DeleteBuilding(ctx, buildingName)
	if err != nil {
		color.Red("Error deleting building: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
	}
}

func (h *CliAdminHandler) AddFloors(ctx context.Context) {
	color.Blue("Enter the name of the building to add floors:")
	var buildingName string
	fmt.Scanln(&buildingName)

	color.Blue("Enter the floor numbers to add (space-separated):")
	var floorNumbersInput string
	fmt.Scanln(&floorNumbersInput)

	var floorNumbers []int
	for _, numStr := range strings.Split(floorNumbersInput, " ") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			color.Red("Invalid floor number: %s", numStr)
			color.Green("Press Enter to continue...")
			fmt.Scanln()
			return
		}
		floorNumbers = append(floorNumbers, num)
	}

	err := h.floorService.AddFloors(ctx, buildingName, floorNumbers)
	if err != nil {
		color.Red("Error adding floors: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
	}
}

func (h *CliAdminHandler) DeleteFloors(ctx context.Context) {
	color.Blue("Enter the name of the building to delete floors from:")
	var buildingName string
	fmt.Scanln(&buildingName)
	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		color.Red("Error fetching floors: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	color.Blue("Available floors in %s: %v", buildingName, floorNumbers)
	color.Blue("Enter the floor numbers to delete(space-separated):")
	var floorNumbersInput string

	fmt.Scanln(&floorNumbersInput)

	var floorsToDelete []int
	for _, numStr := range strings.Split(floorNumbersInput, " ") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			color.Red("Invalid floor number: %s", numStr)
			color.Green("Press Enter to continue...")
			fmt.Scanln()
			return
		}
		floorsToDelete = append(floorsToDelete, num)
	}

	err = h.floorService.DeleteFloors(ctx, buildingName, floorsToDelete)
	if err != nil {
		color.Red("Error deleting floors: %v", err)
	} else {
		color.Green("Floors deleted successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) AddSlots(ctx context.Context) {
	color.Blue("Enter the name of the building to add slots:")
	var buildingName string
	fmt.Scanln(&buildingName)

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		color.Red("Error fetching floors: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	color.Blue("Available floors in %s: %v", buildingName, floorNumbers)

	color.Blue("Enter the floor number to add slots:")
	var floorNumber int
	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		color.Red("Floor %d does not exist in building %s", floorNumber, buildingName)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	color.Blue("Enter the slot numbers to add (space-separated):")
	var slotNumbersInput string
	fmt.Scanln(&slotNumbersInput)

	color.Blue("Enter type of slots (0 - Two Wheeler, 1 - Four Wheeler):")
	var slotType vehicletypes.VehicleType
	fmt.Scanln(&slotType)

	var slotNumbers []int
	for _, numStr := range strings.Split(slotNumbersInput, " ") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			color.Red("Invalid slot number: %s", numStr)
			color.Green("Press Enter to continue...")
			fmt.Scanln()
			return
		}
		slotNumbers = append(slotNumbers, num)
	}
	err = h.slotService.AddSlots(ctx, buildingName, floorNumber, slotNumbers, slotType)

	if err != nil {
		color.Red("Error adding slots: %v", err)
	} else {
		color.Green("Slots added successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteSlots(ctx context.Context) {
	color.Blue("Enter the name of the building to delete slots from:")
	var buildingName string
	fmt.Scanln(&buildingName)

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		color.Red("Error fetching floors: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	color.Blue("Available floors in %s: %v", buildingName, floorNumbers)

	color.Blue("Enter the floor number to delete slots from:")
	var floorNumber int
	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		color.Red("Floor %d does not exist in building %s", floorNumber, buildingName)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	availableSlots, err := h.slotService.GetSlotsByFloor(ctx, buildingName, floorNumber)
	if err != nil {
		color.Red("Error fetching slots: %v", err)
		color.Green("Press Enter to continue...")
		fmt.Scanln()
		return
	}

	color.Blue("Available slots on floor %d in %s: %v", floorNumber, buildingName, availableSlots)

	color.Blue("Enter the slot numbers to delete (space-separated):")
	var slotNumbersInput string
	fmt.Scanln(&slotNumbersInput)

	var slotsToDelete []int
	for _, numStr := range strings.Split(slotNumbersInput, " ") {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			color.Red("Invalid slot number: %s", numStr)
			color.Green("Press Enter to continue...")
			fmt.Scanln()
			return
		}
		slotsToDelete = append(slotsToDelete, num)
	}

	err = h.slotService.DeleteSlots(ctx, buildingName, floorNumber, slotsToDelete)
	if err != nil {
		color.Red("Error deleting slots: %v", err)
	} else {
		color.Green("Slots deleted successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}
