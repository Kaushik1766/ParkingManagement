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
		name             string
		setupMocks       func(*mocks.MockVehicleMgr, *mocks.MockUserManager)
		setupContext     func() context.Context
		setupRequest     func() *http.Request
		expectedStatus   int
		expectedBody     string
		expectedVehicles int // -1 means don't check vehicle count
	}{
		{
			name: "customer gets own vehicles successfully",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().GetRegisteredVehicles(gomock.Any()).Return([]models.VehicleDTO{
					{
						NumberPlate: "ABC123",
						VehicleType: "TwoWheeler",
						AssignedSlot: &models.Slot{
							SlotNumber: 1,
						},
					},
					{
						NumberPlate: "XYZ789",
						VehicleType: "FourWheeler",
						AssignedSlot: &models.Slot{
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
			expectedStatus:   http.StatusOK,
			expectedVehicles: 2,
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
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     `{"message":"Failed to fetch vehicles","code":500}`,
			expectedVehicles: -1,
		},
		{
			name: "admin gets vehicles by userId successfully",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().GetVehiclesByUserId(gomock.Any(), "user-456").Return([]models.VehicleDTO{
					{
						NumberPlate: "DEF456",
						VehicleType: "TwoWheeler",
						AssignedSlot: &models.Slot{
							SlotNumber: 3,
						},
					},
					{
						NumberPlate: "GHI789",
						VehicleType: "FourWheeler",
						AssignedSlot: &models.Slot{
							SlotNumber: 4,
						},
					},
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "admin-123"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/vehicles?userId=user-456", nil)
				return req
			},
			expectedStatus:   http.StatusOK,
			expectedVehicles: 2,
		},
		{
			name: "admin gets vehicles by userId service error",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().GetVehiclesByUserId(gomock.Any(), "user-456").Return(nil, errors.New("user not found"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "admin-123"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/vehicles?userId=user-456", nil)
				return req
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     "{}",
			expectedVehicles: -1,
		},
		{
			name: "admin without userId parameter",
			setupMocks: func(mockVehicleService *mocks.MockVehicleMgr, mockUserService *mocks.MockUserManager) {
				// No service calls expected when userId is not provided
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{Subject: "admin-123"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/vehicles", nil)
			},
			expectedStatus:   http.StatusOK,
			expectedVehicles: 0,
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

			// For successful cases, verify the response structure if expectedVehicles is not -1
			if tt.expectedStatus == http.StatusOK && tt.expectedVehicles != -1 {
				var vehicles []models.VehicleDTO
				if err := json.Unmarshal(w.Body.Bytes(), &vehicles); err != nil {
					t.Errorf("GetVehicles() failed to unmarshal response: %v", err)
				}
				if len(vehicles) != tt.expectedVehicles {
					t.Errorf("GetVehicles() returned %d vehicles, want %d", len(vehicles), tt.expectedVehicles)
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

func TestWebVehicleHandler_RemoveVehicle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockUserManager)
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful vehicle removal",
			setupMock: func(mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().UnregisterVehicle(gomock.Any(), "ABC123").Return(nil)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/vehicles/ABC123", nil)
				req.SetPathValue("numberplate", "ABC123")
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "service error during unregistration",
			setupMock: func(mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().UnregisterVehicle(gomock.Any(), "XYZ789").Return(errors.New("vehicle not found"))
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/vehicles/XYZ789", nil)
				req.SetPathValue("numberplate", "XYZ789")
				return req
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"vehicle not found","code":400}`,
		},
		{
			name: "empty numberplate parameter",
			setupMock: func(mockUserService *mocks.MockUserManager) {
				mockUserService.EXPECT().UnregisterVehicle(gomock.Any(), "").Return(errors.New("invalid numberplate"))
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/vehicles/", nil)
				req.SetPathValue("numberplate", "")
				return req
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid numberplate","code":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVehicleService := mocks.NewMockVehicleMgr(ctrl)
			mockUserService := mocks.NewMockUserManager(ctrl)
			tt.setupMock(mockUserService)

			handler := NewWebVehicleHandler(mockVehicleService, mockUserService)
			w := httptest.NewRecorder()
			req := tt.setupRequest()

			handler.RemoveVehicle(context.Background(), w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("RemoveVehicle() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("RemoveVehicle() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set for error responses
			if tt.expectedBody != "" && w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("RemoveVehicle() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}
		})
	}
}
