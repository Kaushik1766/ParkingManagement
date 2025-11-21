package authservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "kaushik@a.com").MaxTimes(2).Return(models.User{
		UserID:   uuid.Nil,
		Name:     "kaushik",
		Email:    "kaushik@a.com",
		Password: "$2a$12$rviLQTF/MbjKGOC2UiFEdO4RBXkSCZnjGYRKQUa4LJC4CPxGB3nl.",
		Role:     0,
		IsActive: true,
		OfficeID: uuid.Nil,
		Office:   models.Office{},
		Vehicles: []models.Vehicle{},
	}, nil)
	mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "unknown@a.com").Return(models.User{},
		errors.New("user nor found"))

	type fields struct {
		userDb userrepository.UserStorage
	}
	type args struct {
		loginReq models.LoginRequestDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid email",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				loginReq: models.LoginRequestDTO{
					Email:    "kaushik@a.com",
					Password: "a",
				},
			},
			want:    "token",
			wantErr: false,
		},
		{
			name: "unregistered email",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				loginReq: models.LoginRequestDTO{
					Email:    "unknown@a.com",
					Password: "a",
				},
			},
			want:    "token",
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				loginReq: models.LoginRequestDTO{
					Email:    "kaush",
					Password: "a",
				},
			},
			want:    "token",
			wantErr: true,
		},
		{
			name: "wrong password",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				loginReq: models.LoginRequestDTO{
					Email:    "kaushik@a.com",
					Password: "ad",
				},
			},
			want:    "token",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := authservice.NewAuthService(tt.fields.userDb)
			got, err := auth.Login(context.Background(), tt.args.loginReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == "" {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockUserRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		MaxTimes(3).
		Return(nil)

	// Add a separate mock for the error case
	mockUserRepoError := mocks.NewMockUserStorage(ctrl)
	mockUserRepoError.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	type fields struct {
		userDb userrepository.UserStorage
	}
	type args struct {
		registerReq models.RegisterRequestDTO
		role        roles.Role
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid email",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				registerReq: models.RegisterRequestDTO{
					Name:     "kaushik",
					Email:    "kaushik@a.com",
					Office:   "watchguard",
					Password: "123",
				},
				role: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				registerReq: models.RegisterRequestDTO{
					Name:     "kaushik",
					Email:    "kaushik",
					Office:   "watchguard",
					Password: "123",
				},
				role: 0,
			},
			wantErr: true,
		},
		{
			name: "very long password",
			fields: fields{
				userDb: mockUserRepo,
			},
			args: args{
				registerReq: models.RegisterRequestDTO{
					Name:   "kaushik",
					Email:  "kaushik@a.com",
					Office: "watchguard",
					Password: func() string {
						pass := ""
						for range 100 {
							pass += "a"
						}
						return pass
					}(),
				},
				role: 0,
			},
			wantErr: true,
		},
		{
			name: "user creation error",
			fields: fields{
				userDb: mockUserRepoError,
			},
			args: args{
				registerReq: models.RegisterRequestDTO{
					Name:     "kaushik",
					Email:    "kaushik@a.com",
					Office:   "watchguard",
					Password: "123",
				},
				role: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := authservice.NewAuthService(tt.fields.userDb)
			if err := auth.Signup(context.Background(), tt.args.registerReq, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)

	type args struct {
		db userrepository.UserStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful creation",
			args: args{
				db: mockUserRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authservice.NewAuthService(tt.args.db)
			if got == nil {
				t.Errorf("NewAuthService() returned nil")
			}
		})
	}
}
