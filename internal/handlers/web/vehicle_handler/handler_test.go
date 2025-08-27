package vehiclehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestNewWebVehicleHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVehicleService := mocks.NewMockVehicleMgr(ctrl)
	mockUserService := mocks.NewMockUserManager(ctrl)

	tests := []struct {
		name           string
		vehicleService *mocks.MockVehicleMgr
		userService    *mocks.MockUserManager
		want           *WebVehicleHandler
	}{
		{
			name:           "successful creation",
			vehicleService: mockVehicleService,
			userService:    mockUserService,
			want: &WebVehicleHandler{
				vehicleService: mockVehicleService,
				userService:    mockUserService,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWebVehicleHandler(tt.vehicleService, tt.userService)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebVehicleHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebVehicleHandler_GetVehicles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMocks     func(*mocks.MockVehicleMgr, *mocks.MockUserManager)
		setupContext   func() context.Context
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "customer gets own vehicles successfully",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().GetRegisteredVehicles(gomock.Any()).Return([]models.VehicleDTO{
					{
						NumberPlate: "ABC123",
						VehicleType: "TwoWheeler",
						AssignedSlot: models.Slot{
							SlotNumber: 1,
						},
					},
					{
						NumberPlate: "XYZ789",
						VehicleType: "FourWheeler",
						AssignedSlot: models.Slot{
							SlotNumber: 2,
						},
					},
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "user-123"},
					Email:            "kaushik@a.com",
					Role:             roles.Customer,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/vehicles", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "admin gets not implemented",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				// No expectations as admin functionality is not implemented
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "admin-456"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/vehicles", nil)
			},
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name: "customer service error",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().GetRegisteredVehicles(gomock.Any()).Return(nil, errors.New("database error"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "user-123"},
					Email:            "kaushik@a.com",
					Role:             roles.Customer,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/vehicles", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to fetch vehicles","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleService := mocks.NewMockVehicleMgr(ctrl)
			mockUserService := mocks.NewMockUserManager(ctrl)
			tt.setupMocks(mockVehicleService, mockUserService)

			handler := NewWebVehicleHandler(mockVehicleService, mockUserService)
			w := httptest.NewRecorder()
			ctx := tt.setupContext()
			req := tt.setupRequest()

			handler.GetVehicles(ctx, w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetVehicles() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("GetVehicles() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("GetVehicles() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}

			// For successful customer case, verify the response structure
			if tt.expectedStatus == http.StatusOK && tt.expectedBody == "" {
				var vehicles []models.VehicleDTO
				if err := json.Unmarshal(w.Body.Bytes(), &vehicles); err != nil {
					t.Errorf("GetVehicles() failed to unmarshal response: %v", err)
				}
				if len(vehicles) != 2 {
					t.Errorf("GetVehicles() returned %d vehicles, want 2", len(vehicles))
				}
			}
		})
	}
}

func TestWebVehicleHandler_RegisterVehicle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMocks     func(*mocks.MockVehicleMgr, *mocks.MockUserManager)
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful vehicle registration with two wheeler",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().RegisterVehicle(gomock.Any(), "ABC123", vehicletypes.TwoWheeler).Return(nil)
			},
			setupRequest: func() *http.Request {
				vehicleReq := models.AddVehicleDTO{
					NumberPlate: "ABC123",
					VehicleType: 0, // TwoWheeler
				}
				body, _ := json.Marshal(vehicleReq)
				return httptest.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "successful vehicle registration with four wheeler",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().RegisterVehicle(gomock.Any(), "XYZ789", vehicletypes.FourWheeler).Return(nil)
			},
			setupRequest: func() *http.Request {
				vehicleReq := models.AddVehicleDTO{
					NumberPlate: "XYZ789",
					VehicleType: 1, // FourWheeler
				}
				body, _ := json.Marshal(vehicleReq)
				return httptest.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json request body",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				// No expectations as the request should fail before reaching the service
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("POST", "/vehicles", bytes.NewBuffer([]byte("invalid json")))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid request payload","code":400}`,
		},
		{
			name: "invalid vehicle type",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				// No expectations as the request should fail before reaching the service
			},
			setupRequest: func() *http.Request {
				vehicleReq := models.AddVehicleDTO{
					NumberPlate: "ABC123",
					VehicleType: 2, // Invalid type
				}
				body, _ := json.Marshal(vehicleReq)
				return httptest.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid vehicle type","code":400}`,
		},
		{
			name: "service error during registration",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().RegisterVehicle(gomock.Any(), "ABC123", vehicletypes.TwoWheeler).Return(errors.New("vehicle already registered"))
			},
			setupRequest: func() *http.Request {
				vehicleReq := models.AddVehicleDTO{
					NumberPlate: "ABC123",
					VehicleType: 0, // TwoWheeler
				}
				body, _ := json.Marshal(vehicleReq)
				return httptest.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"vehicle already registered","code":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleService := mocks.NewMockVehicleMgr(ctrl)
			mockUserService := mocks.NewMockUserManager(ctrl)
			tt.setupMocks(mockVehicleService, mockUserService)

			handler := NewWebVehicleHandler(mockVehicleService, mockUserService)
			w := httptest.NewRecorder()
			req := tt.setupRequest()

			handler.RegisterVehicle(context.Background(), w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("RegisterVehicle() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("RegisterVehicle() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("RegisterVehicle() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}
		})
	}
}
