package menuconstants

const (
	EnterYourChoice = "Enter your choice:"
	Exit            = "Exit"
	Logout          = "Logout"

	// main menu
	WelcomeMessage    = "Welcome to Parking Management System"
	LoginOption       = "1. Login"
	SignupOption      = "2. Signup"
	AdminSignupOption = "3. Admin Signup(for demo)"
	ExitOption        = "4. Exit"
	LoginRetryMessage = "Please try again or signup if you don't have an account."

	// customer menu
	CustomerUpdateProfile     = "1. Update profile"
	CustomerRegisterVehicle   = "2. Register vehicle"
	CustomerViewProfile       = "3. View Profile"
	CustomerViewVehicles      = "4. View Registered Vehicles"
	CustomerUnregisterVehicle = "5. Unregister Vehicle"
	CustomerParkingMenu       = "6. Parking Menu"
	CustomerLogout            = "7. Logout"
	CustomerExit              = "8. Exit"

	// admin menu
	AdminPageTitle          = "Admin page: "
	AdminBuildingManagement = "1. Building Management"
	AdminFloorManagement    = "2. Floor Management"
	AdminSlotManagement     = "3. Slot Management"
	AdminOfficeManagement   = "4. Office Management"
	AdminUnassignedSlotMgmt = "5. Unassigned Slot Management"
	AdminLogout             = "6. Logout"
	AdminExit               = "7. Exit"

	// parking menu
	ParkingMenuTitle = "Parking Menu:"
	ParkVehicle      = "1. Park Vehicle"
	UnparkVehicle    = "2. Unpark Vehicle"
	ViewParkings     = "3. View Parkings"
	ParkingExit      = "4. Exit"

	// unassigned slot menu
	ViewUnassignedSlots = "1. View Vehicles with Unassigned Slots"
	AssignSlot          = "2. Assign Slot to Vehicle"
	UnassignedSlotExit  = "3. Exit"

	// building management menu
	AddBuilding    = "1. Add Building"
	DeleteBuilding = "2. Delete Building"
	ListBuildings  = "3. List Buildings"
	BuildingExit   = "4. Exit"

	// floor management menu
	AddFloor    = "1. Add Floor"
	DeleteFloor = "2. Delete Floor"
	ListFloors  = "3. List Floors"
	FloorExit   = "4. Exit"

	// slot management menu
	AddSlots    = "1. Add Slots"
	DeleteSlots = "2. Delete Slots"
	ViewSlots   = "3. View Slots"
	SlotExit    = "4. Exit"

	// office management menu
	AddOffice    = "1. Add Office"
	RemoveOffice = "2. Remove Office"
	ListOffices  = "3. List Offices"
	OfficeExit   = "4. Exit"

	// handler common messages
	PressEnterToContinue = "Press Enter to continue..."

	// building operations
	AvailableBuildings         = "Available buildings:"
	EnterBuildingNumber        = "Enter the number of the building to delete:"
	SelectBuildingNumber       = "Select the number of the building :"
	SelectBuildingToDeleteFrom = "Select the number of the building to delete slots from:"
	SelectBuildingToList       = "Select the number of the building to list floors:"
	SelectBuildingToListSlots  = "Select the number of the building to list slots:"
	AddBuildingPrompt          = "Enter name of the building to add: "
	BuildingAddedSuccess       = "Building added successfully."
	BuildingDeletedSuccess     = "Building deleted successfully."

	// floor operations
	AvailableFloors            = "Available floors in %s:"
	SelectBuildingForFloors    = "Enter the number of the building to add floors to:"
	SelectBuildingDeleteFloors = "Enter the number of the building to delete floors from:"
	EnterFloorNumber           = "Enter the floor number to add slots:"
	EnterFloorNumberDelete     = "Enter the floor number to delete slots from:"
	EnterFloorNumberList       = "Enter the floor number to list slots:"
	AddFloorPrompt             = "Enter the floor numbers (dont include already present floors) to add (space-separated):"
	DeleteFloorPrompt          = "Enter space spearated floor numbers to delete (dont enter index numbers):"
	FloorsAddedSuccess         = "Floors added successfully."
	FloorsDeletedSuccess       = "Floors deleted successfully."

	// slot operations
	AvailableSlots           = "Available slots in Floor %d of %s:"
	AvailableSlotsWithStatus = "Available slots in Floor %d of %s (red-> Occupied, green-> vacant):"
	EnterSlotType            = "Enter type of slots (0 - Two Wheeler, 1 - Four Wheeler):"
	AddSlotPrompt            = "Enter the slot numbers to add (space-separated):"
	DeleteSlotPrompt         = "Enter the slot numbers to delete (space-separated):"
	SlotsAddedSuccess        = "Slots added successfully."
	SlotsDeletedSuccess      = "Slots deleted successfully."

	// office operations
	AvailableBuildingsOffice  = "Available Buildings:"
	SelectBuildingForOffice   = "Enter building number to add office:"
	AvailableFloorsOffice     = "Available Floors in %s:"
	SelectFloorForOffice      = "Enter floor number to add office:"
	EnterOfficeName           = "Office Name: "
	OfficeAddedSuccess        = "Office %s added successfully in building %s on floor %d"
	SelectFloorToRemoveOffice = "Enter the floor number of office to remove:"
	OfficeRemovedSuccess      = "Office removed successfully"
	OfficesInBuilding         = "Offices in building %s:"

	// slot assignment operations
	SelectVehicleForSlot = "Select the number of vehicle to assign a slot to:"
	SelectSlotForVehicle = "Select a slot to assign to the vehicle:"
	SlotAssignedSuccess  = "Slot assigned successfully to vehicle %s"

	// parking operations
	SelectVehicleToPark    = "Select the vehicle to park:"
	VehicleParkedSuccess   = "Vehicle parked successfully. Ticket ID: %s"
	SelectParkingToUnpark  = "Select the parking to unpark:"
	VehicleUnparkedSuccess = "Vehicle unparked successfully. Ticket ID: %s"
	ParkingHistoryTitle    = "Parking History:"
)
