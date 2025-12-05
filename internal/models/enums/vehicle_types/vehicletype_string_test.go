package vehicletypes

import "testing"

func TestVehicleType_String(t *testing.T) {
	tests := []struct {
		name string
		i    VehicleType
		want string
	}{
		{
			name: "two wheeler type",
			i:    TwoWheeler,
			want: "TwoWheeler",
		},
		{
			name: "four wheeler type",
			i:    FourWheeler,
			want: "FourWheeler",
		},
		{
			name: "negative vehicle type",
			i:    VehicleType(-1),
			want: "VehicleType(-1)",
		},
		{
			name: "invalid vehicle type value",
			i:    VehicleType(5),
			want: "VehicleType(5)",
		},
		{
			name: "boundary vehicle type",
			i:    VehicleType(2),
			want: "VehicleType(2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
