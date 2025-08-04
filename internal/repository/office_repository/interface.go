package officerepository

type OfficeStorage interface {
	AddOffice(officeName, buildingName string, floorNumber int) error
	DeleteOffice(officeName string) error
}
