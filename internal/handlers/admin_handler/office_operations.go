package adminhandler

import (
	"context"
	"fmt"

	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/fatih/color"
)

func (h *CliAdminHandler) AddOffice(ctx context.Context) {
	color.Yellow("Available Buildings:")
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildings)

	color.Yellow("Enter building number to add office:")
	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow("Available Floors in %s:", buildingName)
	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching floors: %v", err))
		return
	}

	floorNumbersStr := make([]string, len(floorNumbers))
	for i, floor := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", floor)
	}
	utils.PrintListInRows(floorNumbersStr)

	color.Yellow("Enter floor number to add office:")
	var floorNumber int
	fmt.Scanf("%d", &floorNumber)

	color.Yellow("Enter office name:")
	officeName, err := utils.ReadAndSanitizeInput("Office Name: ", &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while reading office name: %v", err))
		return
	}

	err = h.officeService.AddOffice(ctx, officeName, buildingName, floorNumber)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while adding office: %v", err))
		return
	}
	color.Green("Office %s added successfully in building %s on floor %d", officeName, buildingName, floorNumber)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) RemoveOffice(ctx context.Context) {
	color.Yellow("Select building number:")
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildings)

	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow("Select floor number:")
	floorNumbers, err := h.floorService.GetFloorsByBuildingId(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching floors: %v", err))
		return
	}

	floorNumbersStr := make([]string, len(floorNumbers))

	for i, floor := range floorNumbers {
		floorNumbersStr[i] = fmt.Sprintf("Floor %d", floor)
	}

	utils.PrintListInRows(floorNumbersStr)

	var floorNumber int
	fmt.Scanf("%d", &floorNumber)

	color.Yellow("Select office number to remove:")
	officesMap, err := h.officeService.ListOfficesByBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching offices: %v", err))
		return
	}

	floorOffices, ok := officesMap[floorNumber]
	if !ok {
		customerrors.DisplayError(fmt.Sprintf("no offices found on floor %d in building %s", floorNumber, buildingName))
		return
	}

	utils.PrintListInRows(floorOffices)

	color.Yellow("Enter office number to remove:")
	var officeNumber int
	fmt.Scanf("%d", &officeNumber)
	officeName := floorOffices[officeNumber-1]
	err = h.officeService.RemoveOffice(ctx, officeName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while removing office: %v", err))
		return
	}
	color.Green("Office %s removed successfully from building %s on floor %d", officeName, buildingName, floorNumber)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (h *CliAdminHandler) ListOffices(ctx context.Context) {
	color.Yellow("Select building number:")
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildings)

	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow("Offices in building %s:", buildingName)
	officesMap, err := h.officeService.ListOfficesByBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching offices: %v", err))
		return
	}

	for floor, offices := range officesMap {
		color.Green("Floor %d:", floor)
		utils.PrintListInRows(offices)
	}

	color.Green("Press Enter to continue...")
	fmt.Scanln()
}
