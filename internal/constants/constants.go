package constants

type ContextKeys string

const User ContextKeys = "user"

const (
	AdminOffice       string = "ADMIN_OFFICE"
	AdminBuilding     string = "ADMIN_BUILDING"
	SlotLayout        string = "000000000000000111111111111111"
	TestBuilding      string = "TEST_BUILDING"
	TestOffice        string = "TEST_OFFICE"
	BillEmailHeader   string = "Your monthly bill generated on %s is here"
	BillEmailTemplate string = `
<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2>Monthly Parking Bill</h2>
<p>Dear User,</p>
<p>Here is your parking bill details.</p>

<table>
  <tr>
    <th>Date</th>
    <th>Start Time</th>
    <th>End Time</th>
    <th>Vehicle Type</th>
    <th>Number Plate</th>
    <th>Duration (Hours)</th>
    <th>Cost</th>
  </tr>
  %s
</table>

<h3>Total Amount: %.2f</h3>

<p>Thank you for using our parking service.</p>

</body>
</html>
`
	// dynamoDB Keys
	PKBuilding      string = "BUILDING"
	PKUser          string = "USER"
	PKOffice        string = "OFFICE"
	PrefixBuilding  string = "BUILDING#"
	PrefixFloor     string = "FLOOR#"
	PrefixSlot      string = "SLOT#"
	PrefixUser      string = "USER#"
	PrefixDetails   string = "DETAILS#"
	PrefixProfile   string = "PROFILE"
	PrefixFloorInfo string = "FLOORINFO#"

	// error Messages
	ErrInvalidBuildingID        string = "invalid building ID"
	ErrAddingOffice             string = "error adding office"
	ErrNotImplemented           string = "not implemented"
	ErrFetchingOffices          string = "error fetching offices by building"
	ErrFetchingAllOffices       string = "error fetching all offices"
	ErrOfficeNotFound           string = "office not found"
	ErrAddingSlot               string = "error adding slot"
	ErrDeletingSlot             string = "error deleting slot"
	ErrFetchingSlots            string = "error fetching slots"
	ErrFetchingFloors           string = "error fetching floors"
	ErrSavingSlot               string = "error saving slot"
	ErrBuildingNotFound         string = "building not found"
	ErrFetchingUser             string = "error fetching user"
	ErrUserNotFound             string = "user not found"
	ErrFetchingUserList         string = "error fetching user list"
	ErrSavingUser               string = "error saving user"
	ErrOfficeIDRequired         string = "office id is required"
	ErrFetchingOffice           string = "error fetching office"
	ErrSpecifiedOfficeNotFound  string = "specified office not found"
	ErrCheckingExistingUser     string = "error checking existing user"
	ErrUserAlreadyExists        string = "user with given email already exists"
	ErrCreatingUser             string = "error creating user"
	ErrAddingBuilding           string = "error adding building"
	ErrCannotAddAdminBuilding   string = "buildingrepo: cannot add admin building"
	ErrCheckingSlotAvailability string = "error checking slot availability"
	ErrSlotAlreadyOccupied      string = "parkingrepo: vehicle is already parked"
	ErrFetchingUserProfile      string = "error fetching user profile"
	ErrUserProfileNotFound      string = "user profile not found"
	ErrAddingParkingHistory     string = "error adding parking history"
	ErrFindingParkingRecord     string = "error finding parking record"
	ErrParkingRecordNotFound    string = "parking record not found or already unparked"
	ErrUnparkingVehicle         string = "error unparking vehicle"
	ErrFetchingParkingHistory   string = "error fetching parking history"
	ErrNoActiveParkingFound     string = "no active parking found for this numberplate"
	ErrUpdatingParkingRecord    string = "error updating parking record"
	ErrAddingFloor              string = "error adding floor"
	ErrUpdatingBuildingDetails  string = "error updating building details"
	ErrDeletingFloor            string = "error deleting floor"
	ErrAddingVehicle            string = "error adding vehicle"
	ErrVehicleAlreadyExists     string = "vehicle with same numberplate already exists"
	ErrVehicleIsParked          string = "vehicle is parked please unpark it first"
	ErrRemovingVehicle          string = "error removing vehicle"
	ErrFetchingVehicles         string = "error fetching vehicles"
	ErrFetchingVehicle          string = "error fetching vehicle"
	ErrVehicleNotFound          string = "vehicle not found"
	ErrCheckingParkingStatus    string = "error checking parking status"
	ErrSavingVehicle            string = "error saving vehicle"
	ErrSavingBill               string = "error saving bill"
	ErrFetchingBill             string = "error fetching bill"
	ErrBuildingNameExists       string = "building with same name already exists"
	ErrFloorExists              string = "floor already exists"

	PKParking string = "PARKING#"
	PKVehicle string = "VEHICLE#"
	PKBill    string = "BILL#"
)
