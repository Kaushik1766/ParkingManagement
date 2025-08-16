package adminhandler

import (
	"bufio"
	"context"
	"fmt"
	"slices"

	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
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
	buildingName, err := utils.ReadAndSanitizeInput(menuconstants.AddBuildingPrompt, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error reading building name: %v", err))
		return
	}
	err = h.buildingService.AddBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding building: %v", err))
		return
	}
	color.Green(menuconstants.BuildingAddedSuccess)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteBuilding(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildingNames)
	color.Yellow(menuconstants.EnterBuildingNumber)
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
	color.Green(menuconstants.BuildingDeletedSuccess)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) AddFloors(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildingNames)
	color.Yellow(menuconstants.SelectBuildingForFloors)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}
	buildingName := buildingNames[buildingNumber-1]

	color.Blue(menuconstants.AvailableFloors, buildingName)
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

	floorNumbers, err := utils.ReadIntList(menuconstants.AddFloorPrompt, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for floor numbers: %v", err))
		return
	}

	err = h.floorService.AddFloors(ctx, buildingName, floorNumbers)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding floors: %v", err))
		return
	}
	color.Green(menuconstants.FloorsAddedSuccess)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteFloors(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow(menuconstants.SelectBuildingDeleteFloors)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}

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

	floorsToDelete, err := utils.ReadIntList(menuconstants.DeleteFloorPrompt, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for floor numbers: %v", err))
		return
	}

	for _, floor := range floorsToDelete {
		if !slices.Contains(floorNumbers, floor) {
			customerrors.DisplayError(fmt.Sprintf("Floor %d does not exist in building %s", floor, buildingName))
			return
		}
	}

	err = h.floorService.DeleteFloors(ctx, buildingName, floorsToDelete)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error deleting floors: %v", err))
	} else {
		color.Green(menuconstants.FloorsDeletedSuccess)
	}
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) AddSlots(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow(menuconstants.SelectBuildingNumber)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}
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

	color.Blue(menuconstants.EnterFloorNumber)
	var floorNumber int
	fmt.Scanln(&floorNumber)

	if !slices.Contains(floorNumbers, floorNumber) {
		customerrors.DisplayError(fmt.Sprintf("Floor %d does not exist in building %s", floorNumber, buildingName))
		return
	}

	slotNumbers, err := utils.ReadIntList(menuconstants.AddSlotPrompt, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for slot numbers: %v", err))
		return
	}

	color.Blue(menuconstants.EnterSlotType)
	var slotType vehicletypes.VehicleType
	fmt.Scanln(&slotType)

	err = h.slotService.AddSlots(ctx, buildingName, floorNumber, slotNumbers, slotType)

	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error adding slots: %v", err))
	} else {
		color.Green(menuconstants.SlotsAddedSuccess)
	}
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) DeleteSlots(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow(menuconstants.SelectBuildingToDeleteFrom)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	if buildingNumber < 1 || buildingNumber > len(buildingNames) {
		customerrors.DisplayError("Invalid building number selected.")
		return
	}
	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue(menuconstants.AvailableFloors, buildingName)

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(floorNumbersStr)

	color.Blue(menuconstants.EnterFloorNumberDelete)
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

	color.Blue(menuconstants.AvailableSlots, floorNumber, buildingName)

	availableSlotsStr := make([]string, len(availableSlots))
	for i, slot := range availableSlots {
		if slot.SlotType == vehicletypes.TwoWheeler {
			availableSlotsStr[i] = fmt.Sprintf("Slot %d 🏍️", slot.SlotNumber)
		} else {
			availableSlotsStr[i] = fmt.Sprintf("Slot %d 🚗", slot.SlotNumber)
		}
	}

	utils.PrintListInRows(availableSlotsStr)

	slotsToDelete, err := utils.ReadIntList(menuconstants.DeleteSlotPrompt, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Invalid input for slot numbers: %v", err))
		return
	}

	err = h.slotService.DeleteSlots(ctx, buildingName, floorNumber, slotsToDelete)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error deleting slots: %v", err))
		return
	} else {
		color.Green(menuconstants.SlotsDeletedSuccess)
	}
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) ListBuildings(ctx context.Context) {
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	color.Blue(menuconstants.AvailableBuildings)
	utils.PrintListInRows(buildingNames)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) ListFloors(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow(menuconstants.SelectBuildingToList)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)

	buildingName := buildingNames[buildingNumber-1]

	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue(menuconstants.AvailableFloors, buildingName)

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}

	utils.PrintListInRows(floorNumbersStr)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) ListSlots(ctx context.Context) {
	color.Blue(menuconstants.AvailableBuildings)
	buildingNames, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildingNames)

	color.Yellow(menuconstants.SelectBuildingToListSlots)
	var buildingNumber int
	fmt.Scanln(&buildingNumber)
	buildingName := buildingNames[buildingNumber-1]
	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Error fetching floors: %v", err))
		return
	}

	color.Blue(menuconstants.AvailableFloors, buildingName)
	floorNumbersStr := make([]string, len(floorNumbers))
	for i, val := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", val)
	}
	utils.PrintListInRows(floorNumbersStr)
	color.Blue(menuconstants.EnterFloorNumberList)
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

	color.Blue(menuconstants.AvailableSlotsWithStatus, floorNumber, buildingName)

	slotsStr := make([]string, len(slots))

	for i, val := range slots {
		var str string
		if val.SlotType == vehicletypes.TwoWheeler {
			if val.OccupantID != nil {
				str = color.RedString("Slot %d 🏍️\t", val.SlotNumber)
			} else {
				str = color.GreenString("Slot %d 🏍️\t", val.SlotNumber)
			}
		} else {
			if val.OccupantID != nil {
				str = color.RedString("Slot %d 🚗\t", val.SlotNumber)
			} else {
				str = color.GreenString("Slot %d 🚗\t", val.SlotNumber)
			}
		}
		slotsStr[i] = str
	}

	utils.PrintListInRows(slotsStr)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}
