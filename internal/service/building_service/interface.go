package buildingservice

type BuildingMgr interface {
	AddBuilding(name string) error
	DeleteBuilding(name string) error
}
