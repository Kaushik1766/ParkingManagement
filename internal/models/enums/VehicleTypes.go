package enums

type VehicleType int

const (
	TwoWheeler VehicleType = iota
	FourWheeler
)

func (v VehicleType) String() string {
	switch v {
	case TwoWheeler:
		return "TwoWheeler"
	default:
		return "FourWheeler"
	}
}
