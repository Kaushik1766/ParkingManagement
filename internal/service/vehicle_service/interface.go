package vehicleservice

import "context"

type VehicleMgr interface {
	Park(ctx context.Context, numberplate string) (string, error)
	Unpark(ctx context.Context, ticketId string) error
}
