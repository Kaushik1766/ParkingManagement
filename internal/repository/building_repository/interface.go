package buildingrepository

type BuildingStorage interface {
	AddBuilding(name string) error
	DeleteBuilding(name string) error
}
