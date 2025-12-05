package roles

import "testing"

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		i    Role
		want string
	}{
		{
			name: "customer role string",
			i:    Customer,
			want: "Customer",
		},
		{
			name: "admin role string",
			i:    Admin,
			want: "Admin",
		},
		{
			name: "invalid negative role",
			i:    Role(-1),
			want: "Role(-1)",
		},
		{
			name: "invalid high role value",
			i:    Role(99),
			want: "Role(99)",
		},
		{
			name: "out of bounds role",
			i:    Role(2),
			want: "Role(2)",
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
