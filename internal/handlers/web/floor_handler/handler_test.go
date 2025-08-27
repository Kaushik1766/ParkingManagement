package floorhandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
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

func TestNewWebFloorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFloorService := mocks.NewMockFloorMgr(ctrl)

	type args struct {
		floorService floorservice.FloorMgr
	}
	tests := []struct {
		name string
		args args
		want *WebFloorHandler
	}{
		{
			name: "valid floor service",
			args: args{
				floorService: mockFloorService,
			},
			want: &WebFloorHandler{
				floorService: mockFloorService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebFloorHandler(tt.args.floorService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebFloorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebFloorHandler_GetFloors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFloorService := mocks.NewMockFloorMgr(ctrl)

	type fields struct {
		floorService floorservice.FloorMgr
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
			name: "admin get floors success",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/floors", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				floors := []models.FloorDTO{
					{BuildingID: "123", FloorNumber: 1},
					{BuildingID: "123", FloorNumber: 2},
				}
				mockFloorService.EXPECT().GetFloorsByBuildingId(gomock.Any(), "123").Return(floors, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"building_id\":\"123\",\"floor_number\":1},{\"building_id\":\"123\",\"floor_number\":2}]\n",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/floors", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusUnauthorized,
				body:   "{\"message\":\"Unauthorized access\",\"code\":401}\n",
			},
		},
		{
			name: "get floors service error",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/floors", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockFloorService.EXPECT().GetFloorsByBuildingId(gomock.Any(), "123").Return(nil, errors.New("building not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"building not found\",\"code\":500}\n",
			},
		},
		{
			name: "empty floors list",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/456/floors", nil)
					req.SetPathValue("buildingId", "456")
					return req
				}(),
			},
			mock: func() {
				floors := []models.FloorDTO{}
				mockFloorService.EXPECT().GetFloorsByBuildingId(gomock.Any(), "456").Return(floors, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[]\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebFloorHandler{
				floorService: tt.fields.floorService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.GetFloors(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("GetFloors() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("GetFloors() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebFloorHandler_AddFloor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFloorService := mocks.NewMockFloorMgr(ctrl)

	type fields struct {
		floorService floorservice.FloorMgr
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
			name: "admin add floor success",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/floors",
						bytes.NewBufferString(`{"floor_number":1}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockFloorService.EXPECT().AddFloorByBuildingId(gomock.Any(), "123", 1).Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/floors",
						bytes.NewBufferString(`{"floor_number":1}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusUnauthorized,
				body:   "{\"message\":\"Unauthorized access\",\"code\":401}\n",
			},
		},
		{
			name: "invalid json request",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/floors",
						bytes.NewBufferString(`{"floor_number":}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request body\",\"code\":400}\n",
			},
		},
		{
			name: "add floor service error",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/floors",
						bytes.NewBufferString(`{"floor_number":2}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockFloorService.EXPECT().AddFloorByBuildingId(gomock.Any(), "123", 2).Return(errors.New("floor already exists"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"floor already exists\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebFloorHandler{
				floorService: tt.fields.floorService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.AddFloor(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("AddFloor() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("AddFloor() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebFloorHandler_DeleteFloor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFloorService := mocks.NewMockFloorMgr(ctrl)

	type fields struct {
		floorService floorservice.FloorMgr
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
			name: "admin delete floor success",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/123/floors/1", nil)
					req.SetPathValue("buildingId", "123")
					req.SetPathValue("floorId", "1")
					return req
				}(),
			},
			mock: func() {
				mockFloorService.EXPECT().DeleteFloor(gomock.Any(), "123", 1).Return(nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/123/floors/1", nil)
					req.SetPathValue("buildingId", "123")
					req.SetPathValue("floorId", "1")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusUnauthorized,
				body:   "{\"message\":\"Unauthorized access\",\"code\":401}\n",
			},
		},
		{
			name: "invalid floor id",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/123/floors/invalid", nil)
					req.SetPathValue("buildingId", "123")
					req.SetPathValue("floorId", "invalid")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid floor ID\",\"code\":400}\n",
			},
		},
		{
			name: "delete floor service error",
			fields: fields{
				floorService: mockFloorService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/buildings/123/floors/2", nil)
					req.SetPathValue("buildingId", "123")
					req.SetPathValue("floorId", "2")
					return req
				}(),
			},
			mock: func() {
				mockFloorService.EXPECT().DeleteFloor(gomock.Any(), "123", 2).Return(errors.New("floor not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"floor not found\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebFloorHandler{
				floorService: tt.fields.floorService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.DeleteFloor(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("DeleteFloor() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("DeleteFloor() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}
