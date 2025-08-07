package parkinghistory

type ParkingHistoryDTO struct {
	TicketId    string
	NumberPlate string
	BuildingId  string
	FLoorNumber int
	SlotNumber  int
	StartTime   string
	EndTime     string
}
