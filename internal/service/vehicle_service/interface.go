package vehicleservice

import "context"

//go:generate mockgen -source=interface.go -destination=../../../mocks/vehicle_service_mock.go -package=mocks
type VehicleMgr interface {
	Park(ctx context.Context, numberplate string) (string, error)
	Unpark(ctx context.Context, ticketId string) error
	UnparkByNumberPlate(ctx context.Context, numberplate string) error
}
