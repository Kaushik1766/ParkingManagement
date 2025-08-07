package adminhandler

import (
	"bufio"
	"context"
	"fmt"
	"slices"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/fatih/color"
)

type CliAdminHandler struct {
	floorService    floorservice.FloorMgr
	buildingService buildingservice.BuildingMgr
	slotService     slotservice.SlotMgr
	officeService   officeservice.OfficeMgr
	reader          bufio.Reader
}

func NewCliAdminHandler(
	floorService floorservice.FloorMgr,
	buildingService buildingservice.BuildingMgr,
	slotService slotservice.SlotMgr,
	reader *bufio.Reader,
	officeService officeservice.OfficeMgr,
) *CliAdminHandler {
	return &CliAdminHandler{
		floorService:    floorService,
		buildingService: buildingService,
		slotService:     slotService,
		reader:          *reader,
		officeService:   officeService,
	}
}

func (h *CliAdminHandler) AddBuilding(ctx context.Context) {
	buildingName, err := utils.ReadAndSanitizeInput("Enter name of the building to add: ", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error reading building name: %v", err))
		return
	}
	err = h.buildingService.AddBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding building: %v", err))
		return
	}
	color.Green("Building added successfully.")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteBuilding(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildingNames)
	color.Yellow("Enter the number of the building to delete:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}
	buildingName := buildingNames[buildingNumber-1]

	err = h.buildingService.DeleteBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error deleting building: %v", err))
		return
	}
	color.Green("Building deleted successfully.")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) AddFloors(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildingNames)
	color.Yellow("Enter the number of the building to add floors to:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}
	buildingName := buildingNames[buildingNumber-1]

	color.Blue("Available floors in %s:", buildingName)
	availableFloorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	availableFloorsStr := make([]string, len(availableFloorNumbers))
	for i, val := range availableFloorNumbers {
		availableFloorsStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(availableFloorsStr)

	floorNumbers, err := utils.ReadIntList("Enter the floor numbers (dont include already present floors) to add (space-separated):", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for floor numbers: %v", err))
		return
	}

	err = h.floorService.AddFloors(ctx, buildingName, floorNumbers)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding floors: %v", err))
		return
	}
	color.Green("Floors added successfully.")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteFloors(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow("Enter the number of the building to delete floors from:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	availableFloorsStr := make([]string, len(floorNumbers))

	for i, val := range floorNumbers {
		availableFloorsStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(availableFloorsStr)

	floorsToDelete, err := utils.ReadIntList("Enter space spearated floor numbers to delete (dont enter index numbers):", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for floor numbers: %v", err))
		return
	}

	err = h.floorService.DeleteFloors(ctx, buildingName, floorsToDelete)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error deleting floors: %v", err))
	} else {
		color.Green("Floors deleted successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) AddSlots(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow("Select the number of the building to add slots:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(floorNumbersStr)

	color.Blue("Enter the floor number to add slots:")
	var floorNumber int
	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		customerrors.DisplayError(fmt.Sprintf("Floor %d does not exist in building %s", floorNumber, buildingName))
		return
	}

	slotNumbers, err := utils.ReadIntList("Enter the slot numbers to add (space-separated):", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for slot numbers: %v", err))
		return
	}

	color.Blue("Enter type of slots (0 - Two Wheeler, 1 - Four Wheeler):")
	var slotType vehicletypes.VehicleType
	fmt.Scanln(&slotType)

	err = h.slotService.AddSlots(ctx, buildingName, floorNumber, slotNumbers, slotType)

	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding slots: %v", err))
	} else {
		color.Green("Slots added successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteSlots(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow("Select the number of the building to delete slots from:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue("Available floors in %s:", buildingName)

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(floorNumbersStr)

	color.Blue("Enter the floor number to delete slots from:")
	var floorNumber int
	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		customerrors.DisplayError(fmt.Sprintf("Floor %d does not exist in building %s", floorNumber, buildingName))
		return
	}

	availableSlots, err := h.slotService.GetSlotsByFloor(ctx, buildingName, floorNumber)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching slots: %v", err))
		return
	}

	color.Blue("Available slots in Floor %d of %s:", floorNumber, buildingName)

	availableSlotsStr := make([]string, len(availableSlots))
	for i, slot := range availableSlots {
		if slot.SlotType == vehicletypes.TwoWheeler {
			availableSlotsStr[i] = fmt.Sprintf("Slot %d 🏍️", slot.SlotNumber)
		} else {
			availableSlotsStr[i] = fmt.Sprintf("Slot %d 🚗", slot.SlotNumber)
		}
	}

	utils.PrintListInRows(availableSlotsStr)

	slotsToDelete, err := utils.ReadIntList("Enter the slot numbers to delete (space-separated):", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for slot numbers: %v", err))
		return
	}

	err = h.slotService.DeleteSlots(ctx, buildingName, floorNumber, slotsToDelete)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error deleting slots: %v", err))
		return
	} else {
		color.Green("Slots deleted successfully.")
	}
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) ListBuildings(ctx context.Context) {
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	color.Blue("Available buildings:")
	utils.PrintListInRows(buildingNames)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) ListFloors(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow("Select the number of the building to list floors:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue("Available floors in %s:", buildingName)

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(floorNumbersStr)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) ListSlots(ctx context.Context) {
	color.Blue("Available buildings:")
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow("Select the number of the building to list slots:")
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	buildingName := buildingNames[buildingNumber-1]
	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue("Available floors in %s:", buildingName)
	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}
	utils.PrintListInRows(floorNumbersStr)
	color.Blue("Enter the floor number to list slots:")
	var floorNumber int

	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		customerrors.DisplayError(fmt.Sprintf("Floor %d does not exist in building %s", floorNumber, buildingName))
		return
	}

	slots, err := h.slotService.GetSlotsByFloor(ctx, buildingName, floorNumber)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching slots: %v", err))
		return
	}

	color.Blue("Available slots in Floor %d of %s (red-> Occupied, green-> vacant):", floorNumber, buildingName)

	slotsStr := make([]string, len(slots))

	for i, val := range slots {
		var str string
		if val.SlotType == vehicletypes.TwoWheeler {
			if val.IsOccupied {
				str = color.RedString("Slot %d 🏍️\t", val.SlotNumber)
			} else {
				str = color.GreenString("Slot %d 🏍️\t", val.SlotNumber)
			}
		} else {
			if val.IsOccupied {
				str = color.RedString("Slot %d 🚗\t", val.SlotNumber)
			} else {
				str = color.GreenString("Slot %d 🚗\t", val.SlotNumber)
			}
		}
		slotsStr[i] = str
	}

	utils.PrintListInRows(slotsStr)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}
