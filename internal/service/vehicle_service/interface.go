package vehicleservice

import "context"

type VehicleMgr interface {
	Park(ctx context.Context, vehicleId string) error
}
