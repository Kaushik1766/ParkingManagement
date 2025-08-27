package buildinghandler

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
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
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

func TestNewWebBuildingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuildingService := mocks.NewMockBuildingMgr(ctrl)

	type args struct {
		buildingService buildingservice.BuildingMgr
	}
	tests := []struct {
		name string
		args args
		want *WebBuildingHandler
	}{
		{
			name: "valid building service",
			args: args{
				buildingService: mockBuildingService,
			},
			want: &WebBuildingHandler{
				buildingService: mockBuildingService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebBuildingHandler(tt.args.buildingService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebBuildingHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebBuildingHandler_AddBuilding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuildingService := mocks.NewMockBuildingMgr(ctrl)

	type fields struct {
		buildingService buildingservice.BuildingMgr
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
			name: "valid admin add building",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/buildings",
					bytes.NewBufferString(`{"building_name":"kaushik building"}`)),
			},
			mock: func() {
				mockBuildingService.EXPECT().AddBuilding(gomock.Any(), "kaushik building").Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "{\"message\":\"Building added successfully\"}\n",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/buildings",
					bytes.NewBufferString(`{"building_name":"kaushik building"}`)),
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
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/buildings",
					bytes.NewBufferString(`{"building_name":}`)),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request payload\",\"code\":400}\n",
			},
		},
		{
			name: "empty building name",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/buildings",
					bytes.NewBufferString(`{"building_name":""}`)),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request payload\",\"code\":400}\n",
			},
		},
		{
			name: "service error",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/buildings",
					bytes.NewBufferString(`{"building_name":"kaushik building"}`)),
			},
			mock: func() {
				mockBuildingService.EXPECT().AddBuilding(gomock.Any(), "kaushik building").Return(errors.New("building already exists"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to add building\",\"code\":500}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := WebBuildingHandler{
				buildingService: tt.fields.buildingService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.AddBuilding(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("AddBuilding() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("AddBuilding() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebBuildingHandler_GetBuildings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuildingService := mocks.NewMockBuildingMgr(ctrl)

	type fields struct {
		buildingService buildingservice.BuildingMgr
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
			name: "admin get all buildings success",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings", nil),
			},
			mock: func() {
				buildings := []models.BuildingDTO{
					{BuildingID: "123", Name: "kaushik building"},
					{BuildingID: "456", Name: "admin building"},
				}
				mockBuildingService.EXPECT().GetAllBuildings(gomock.Any()).Return(buildings, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"building_id\":\"123\",\"name\":\"kaushik building\"},{\"building_id\":\"456\",\"name\":\"admin building\"}]\n",
			},
		},
		{
			name: "admin get building by id success",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings?buildingId=123", nil),
			},
			mock: func() {
				building := models.BuildingDTO{BuildingID: "123", Name: "kaushik building"}
				mockBuildingService.EXPECT().GetBuildingByID(gomock.Any(), "123").Return(building, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"building_id\":\"123\",\"name\":\"kaushik building\"}]\n",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings", nil),
			},
			mock: func() {},
			want: want{
				status: http.StatusUnauthorized,
				body:   "{\"message\":\"Unauthorized access\",\"code\":401}\n",
			},
		},
		{
			name: "get building by id error",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings?buildingId=123", nil),
			},
			mock: func() {
				mockBuildingService.EXPECT().GetBuildingByID(gomock.Any(), "123").Return(models.BuildingDTO{}, errors.New("building not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to fetch building\",\"code\":500}\n",
			},
		},
		{
			name: "get all buildings error",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings", nil),
			},
			mock: func() {
				mockBuildingService.EXPECT().GetAllBuildings(gomock.Any()).Return(nil, errors.New("database error"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to fetch buildings\",\"code\":500}\n",
			},
		},
		{
			name: "empty buildings list",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequest(http.MethodGet, "/buildings", nil),
			},
			mock: func() {
				buildings := []models.BuildingDTO{}
				mockBuildingService.EXPECT().GetAllBuildings(gomock.Any()).Return(buildings, nil)
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
			handler := WebBuildingHandler{
				buildingService: tt.fields.buildingService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.GetBuildings(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("GetBuildings() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("GetBuildings() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebBuildingHandler_DeleteBuilding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuildingService := mocks.NewMockBuildingMgr(ctrl)

	type fields struct {
		buildingService buildingservice.BuildingMgr
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
			name: "admin delete building success",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/buildings/123", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockBuildingService.EXPECT().DeleteBuildingByID(gomock.Any(), "123").Return(nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "",
			},
		},
		{
			name: "unauthorized customer access",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/buildings/123", nil)
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
			name: "delete building service error",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/buildings/123", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockBuildingService.EXPECT().DeleteBuildingByID(gomock.Any(), "123").Return(errors.New("building not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to delete building\",\"code\":500}\n",
			},
		},
		{
			name: "delete admin building",
			fields: fields{
				buildingService: mockBuildingService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/buildings/456", nil)
					req.SetPathValue("buildingId", "456")
					return req
				}(),
			},
			mock: func() {
				mockBuildingService.EXPECT().DeleteBuildingByID(gomock.Any(), "456").Return(nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebBuildingHandler{
				buildingService: tt.fields.buildingService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.DeleteBuilding(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("DeleteBuilding() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("DeleteBuilding() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}
