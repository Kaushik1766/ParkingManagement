package slothandler

import (
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
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"go.uber.org/mock/gomock"
)

func TestNewWebSlotHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSlotService := mocks.NewMockSlotMgr(ctrl)

	tests := []struct {
		name        string
		slotService *mocks.MockSlotMgr
		want        *WebSlotHandler
	}{
		{
			name:        "successful creation",
			slotService: mockSlotService,
			want: &WebSlotHandler{
				slotService: mockSlotService,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWebSlotHandler(tt.slotService)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebSlotHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebSlotHandler_GetSlots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockSlotMgr)
		setupContext   func() context.Context
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get slots for admin user",
			setupMock: func(mockService *mocks.MockSlotMgr) {
				mockService.EXPECT().GetSlotsByFloor(
					gomock.Any(),
					"building-123",
					1,
				).Return([]models.SlotDTO{
					{
						BuildingID:  "building-123",
						FloorNumber: 1,
						SlotNumber:  1,
						SlotType:    "car",
						IsOccupied:  false,
					},
					{
						BuildingID:  "building-123",
						FloorNumber: 1,
						SlotNumber:  2,
						SlotType:    "bike",
						IsOccupied:  true,
					},
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					Email:  "admin@a.com",
					Role:   roles.Admin,
					Office: "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/slots", nil)
				req.SetPathValue("buildingId", "building-123")
				req.SetPathValue("floorId", "1")
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized access for customer user",
			setupMock: func(mockService *mocks.MockSlotMgr) {
				// No expectations as the request should fail before reaching the service
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					Email:  "kaushik@a.com",
					Role:   roles.Customer,
					Office: "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/slots", nil)
				req.SetPathValue("buildingId", "building-123")
				req.SetPathValue("floorId", "1")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Unauthorized access","code":401}`,
		},
		{
			name: "invalid floor id",
			setupMock: func(mockService *mocks.MockSlotMgr) {
				// No expectations as the request should fail before reaching the service
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					Email:  "admin@a.com",
					Role:   roles.Admin,
					Office: "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/slots", nil)
				req.SetPathValue("buildingId", "building-123")
				req.SetPathValue("floorId", "invalid")
				return req
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid floor ID","code":400}`,
		},
		{
			name: "service error",
			setupMock: func(mockService *mocks.MockSlotMgr) {
				mockService.EXPECT().GetSlotsByFloor(
					gomock.Any(),
					"building-123",
					1,
				).Return(nil, errors.New("database error"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					Email:  "admin@a.com",
					Role:   roles.Admin,
					Office: "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/slots", nil)
				req.SetPathValue("buildingId", "building-123")
				req.SetPathValue("floorId", "1")
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to fetch slots","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockSlotMgr(ctrl)
			tt.setupMock(mockService)

			handler := NewWebSlotHandler(mockService)
			w := httptest.NewRecorder()
			ctx := tt.setupContext()
			req := tt.setupRequest()

			handler.GetSlots(ctx, w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetSlots() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("GetSlots() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("GetSlots() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}

			if tt.expectedStatus == http.StatusOK && tt.expectedBody == "" {
				var slots []models.SlotDTO
				if err := json.Unmarshal(w.Body.Bytes(), &slots); err != nil {
					t.Errorf("GetSlots() failed to unmarshal response: %v", err)
				}
				if len(slots) != 2 {
					t.Errorf("GetSlots() returned %d slots, want 2", len(slots))
				}
			}
		})
	}
}
