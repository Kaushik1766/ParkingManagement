package adminhandler

import (
	"context"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/fatih/color"
)

func (h *CliAdminHandler) AddOffice(ctx context.Context) {
	color.Yellow(menuconstants.AvailableBuildingsOffice)
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}

	utils.PrintListInRows(buildings)

	color.Yellow(menuconstants.SelectBuildingForOffice)
	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow(menuconstants.AvailableFloorsOffice, buildingName)
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

	color.Yellow(menuconstants.SelectFloorForOffice)
	var floorNumber int
	fmt.Scanf("%d", &floorNumber)

	color.Yellow("Enter office name:")
	officeName, err := utils.ReadAndSanitizeInput(menuconstants.EnterOfficeName, &h.reader)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while reading office name: %v", err))
		return
	}

	err = h.officeService.AddOffice(ctx, officeName, buildingName, floorNumber)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while adding office: %v", err))
		return
	}
	color.Green(menuconstants.OfficeAddedSuccess, officeName, buildingName, floorNumber)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) RemoveOffice(ctx context.Context) {
	color.Yellow(menuconstants.SelectBuildingNumber)
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildings)

	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow(menuconstants.SelectFloorToRemoveOffice)
	offices, err := h.officeService.ListOfficesByBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching offices: %v", err))
		return
	}

	for floor, office := range offices {
		color.Green("Floor %d: %s", floor, office)
	}
	var floorNumber int
	fmt.Scanf("%d", &floorNumber)

	officeName := offices[floorNumber]

	err = h.officeService.RemoveOffice(ctx, officeName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while removing office: %v", err))
		return
	}
	color.Green(menuconstants.OfficeRemovedSuccess)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (h *CliAdminHandler) ListOffices(ctx context.Context) {
	color.Yellow(menuconstants.SelectBuildingNumber)
	buildings, err := h.buildingService.GetAllBuildings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching buildings: %v", err))
		return
	}
	utils.PrintListInRows(buildings)

	var buildingNumber int
	fmt.Scanf("%d", &buildingNumber)

	buildingName := buildings[buildingNumber-1]

	color.Yellow(menuconstants.OfficesInBuilding, buildingName)
	officesMap, err := h.officeService.ListOfficesByBuilding(ctx, buildingName)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("error while fetching offices: %v", err))
		return
	}

	for floor, office := range officesMap {
		color.Green("Floor %d: %v", floor, office)
	}

	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}
