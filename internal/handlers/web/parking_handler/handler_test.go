package parkinghandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"go.uber.org/mock/gomock"
)

var (
	adminCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
		Email:  "admin@a.com",
		Role:   roles.Admin,
		Office: "wg",
	})

	customerCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
		Email:  "kaushik@a.com",
		Role:   roles.Customer,
		Office: "wg",
	})
)

func TestNewWebParkingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingService := mocks.NewMockParkingHistoryMgr(ctrl)
	mockVehicleService := mocks.NewMockVehicleMgr(ctrl)

	type args struct {
		parkingService parkinghistoryservice.ParkingHistoryMgr
		vehicleService vehicleservice.VehicleMgr
	}
	tests := []struct {
		name string
		args args
		want *WebParkingHandler
	}{
		{
			name: "valid services",
			args: args{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			want: &WebParkingHandler{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebParkingHandler(tt.args.parkingService, tt.args.vehicleService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebParkingHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebParkingHandler_GetParkings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingService := mocks.NewMockParkingHistoryMgr(ctrl)
	mockVehicleService := mocks.NewMockVehicleMgr(ctrl)

	type fields struct {
		parkingService parkinghistoryservice.ParkingHistoryMgr
		vehicleService vehicleservice.VehicleMgr
	}
	type args struct {
		ctx context.Context
		w   http.ResponseWriter
		r   *http.Request
	}
	type want struct {
		status int
		body   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
		want   want
	}{
		{
			name: "get parkings success with default times",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/api/v1/parkings", nil),
			},
			mock: func() {
				parkings := []models.ParkingHistoryDTO{
					{
						TicketId:     "ticket123",
						NumberPlate:  "kaushik123",
						BuildingId:   "building1",
						FLoorNumber:  1,
						SlotNumber:   5,
						StartTime:    time.Date(2025, 8, 27, 10, 0, 0, 0, time.UTC),
						EndTime:      time.Date(2025, 8, 27, 12, 0, 0, 0, time.UTC),
						VechicleType: vehicletypes.FourWheeler,
					},
				}
				mockParkingService.EXPECT().GetParkingHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(parkings, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"TicketId\":\"ticket123\",\"NumberPlate\":\"kaushik123\",\"BuildingId\":\"building1\",\"FLoorNumber\":1,\"SlotNumber\":5,\"StartTime\":\"2025-08-27T10:00:00Z\",\"EndTime\":\"2025-08-27T12:00:00Z\",\"VechicleType\":1}]\n",
			},
		},
		{
			name: "get parkings with custom times",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/api/v1/parkings?startTime=2025-08-01T00:00:00Z&endTime=2025-08-28T00:00:00Z", nil),
			},
			mock: func() {
				parkings := []models.ParkingHistoryDTO{}
				mockParkingService.EXPECT().GetParkingHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(parkings, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[]\n",
			},
		},
		{
			name: "invalid start time format",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/api/v1/parkings?startTime=invalid", nil),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid startTime format. Use RFC3339 format.\",\"code\":400}\n",
			},
		},
		{
			name: "invalid end time format",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/api/v1/parkings?endTime=invalid", nil),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid endTime format. Use RFC3339 format.\",\"code\":400}\n",
			},
		},
		{
			name: "get parkings service error",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/api/v1/parkings", nil),
			},
			mock: func() {
				mockParkingService.EXPECT().GetParkingHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to fetch parking history: database error\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebParkingHandler{
				parkingService: tt.fields.parkingService,
				vehicleService: tt.fields.vehicleService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.GetParkings(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("GetParkings() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("GetParkings() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebParkingHandler_AddParking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingService := mocks.NewMockParkingHistoryMgr(ctrl)
	mockVehicleService := mocks.NewMockVehicleMgr(ctrl)

	type fields struct {
		parkingService parkinghistoryservice.ParkingHistoryMgr
		vehicleService vehicleservice.VehicleMgr
	}
	type args struct {
		ctx context.Context
		w   http.ResponseWriter
		r   *http.Request
	}
	type want struct {
		status int
		body   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
		want   want
	}{
		{
			name: "add parking success",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/v1/parkings",
					bytes.NewBufferString(`{"numberplate":"kaushik123"}`)),
			},
			mock: func() {
				mockVehicleService.EXPECT().Park(gomock.Any(), "kaushik123").Return("ticket456", nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "{\"ticketId\":\"ticket456\"}\n",
			},
		},
		{
			name: "customer add parking success",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/v1/parkings",
					bytes.NewBufferString(`{"numberplate":"admin789"}`)),
			},
			mock: func() {
				mockVehicleService.EXPECT().Park(gomock.Any(), "admin789").Return("ticket789", nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "{\"ticketId\":\"ticket789\"}\n",
			},
		},
		{
			name: "invalid json request",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/v1/parkings",
					bytes.NewBufferString(`{"numberplate":}`)),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request payload: invalid character '}' looking for beginning of value\",\"code\":400}\n",
			},
		},
		{
			name: "park vehicle service error",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/v1/parkings",
					bytes.NewBufferString(`{"numberplate":"invalid123"}`)),
			},
			mock: func() {
				mockVehicleService.EXPECT().Park(gomock.Any(), "invalid123").Return("", errors.New("no slots available"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to park vehicle: no slots available\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebParkingHandler{
				parkingService: tt.fields.parkingService,
				vehicleService: tt.fields.vehicleService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.AddParking(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("AddParking() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("AddParking() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebParkingHandler_UnparkVehicle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingService := mocks.NewMockParkingHistoryMgr(ctrl)
	mockVehicleService := mocks.NewMockVehicleMgr(ctrl)

	type fields struct {
		parkingService parkinghistoryservice.ParkingHistoryMgr
		vehicleService vehicleservice.VehicleMgr
	}
	type args struct {
		ctx context.Context
		w   http.ResponseWriter
		r   *http.Request
	}
	type want struct {
		status int
		body   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
		want   want
	}{
		{
			name: "unpark vehicle success",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/parkings/kaushik123", nil)
					req.SetPathValue("numberplate", "kaushik123")
					return req
				}(),
			},
			mock: func() {
				mockVehicleService.EXPECT().UnparkByNumberPlate(gomock.Any(), "kaushik123").Return(nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "customer unpark vehicle success",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/parkings/admin789", nil)
					req.SetPathValue("numberplate", "admin789")
					return req
				}(),
			},
			mock: func() {
				mockVehicleService.EXPECT().UnparkByNumberPlate(gomock.Any(), "admin789").Return(nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "unpark vehicle service error",
			fields: fields{
				parkingService: mockParkingService,
				vehicleService: mockVehicleService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/parkings/notfound123", nil)
					req.SetPathValue("numberplate", "notfound123")
					return req
				}(),
			},
			mock: func() {
				mockVehicleService.EXPECT().UnparkByNumberPlate(gomock.Any(), "notfound123").Return(errors.New("vehicle not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to unpark vehicle: vehicle not found\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebParkingHandler{
				parkingService: tt.fields.parkingService,
				vehicleService: tt.fields.vehicleService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.UnparkVehicle(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("UnparkVehicle() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("UnparkVehicle() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}
