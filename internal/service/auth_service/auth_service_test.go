package authservice_test

import (
	"testing"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
)

func TestAuthService_Signup(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db       userrepository.UserStorage
		officeDb officerepository.OfficeStorage
		// Named input parameters for target function.
		registerReq models.RegisterRequestDTO
		role        roles.Role
		wantErr     bool
	}{
		{
			name: "valid data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := authservice.NewAuthService(tt.db)
			gotErr := auth.Signup(tt.registerReq, tt.role)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Signup() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Signup() succeeded unexpectedly")
			}
		})
	}
}
