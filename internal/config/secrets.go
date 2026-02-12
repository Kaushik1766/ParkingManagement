package config

import "time"

const (
	JWTSecret          = "asdfasasdfasdf"
	BillingDuration    = time.Hour * 1
	UsersPath          = "data/users.json"
	BuildingsPath      = "data/buildings.json"
	FloorsPath         = "data/floors.json"
	OfficesPath        = "data/offices.json"
	ParkingHistoryPath = "data/parking_history.json"
	SlotsPath          = "data/slots.json"
	VehiclesPath       = "data/vehicles.json"
	EnvPath            = "./.env"
	DynamoDBTable      = "parking-management-test"
	AwsRegion          = "ap-south-1"
)
