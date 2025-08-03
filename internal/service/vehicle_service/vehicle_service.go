package vehicleservice

import "context"

type VehicleService struct{}

func (vs *VehicleService) Park(ctx context.Context, numberplate string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (vs *VehicleService) Unpark(ctx context.Context, numberplate string) error {
	panic("not implemented") // TODO: Implement
}
