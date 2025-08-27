package officehandler

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
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
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

func TestNewWebOfficeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeService := mocks.NewMockOfficeMgr(ctrl)

	type args struct {
		officeService officeservice.OfficeMgr
	}
	tests := []struct {
		name string
		args args
		want *WebOfficeHandler
	}{
		{
			name: "valid office service",
			args: args{
				officeService: mockOfficeService,
			},
			want: &WebOfficeHandler{
				officeService: mockOfficeService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebOfficeHandler(tt.args.officeService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebOfficeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebOfficeHandler_GetOffices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeService := mocks.NewMockOfficeMgr(ctrl)

	type fields struct {
		officeService officeservice.OfficeMgr
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
			name: "get offices success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/offices", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				offices := []models.OfficeDTO{
					{BuildingID: "123", FloorNumber: 1, OfficeName: "kaushik office", OfficeID: "office1"},
					{BuildingID: "123", FloorNumber: 2, OfficeName: "admin office", OfficeID: "office2"},
				}
				mockOfficeService.EXPECT().ListOfficesByBuilding(gomock.Any(), "123").Return(offices, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"building_id\":\"123\",\"floor_number\":1,\"office_name\":\"kaushik office\",\"office_id\":\"office1\"},{\"building_id\":\"123\",\"floor_number\":2,\"office_name\":\"admin office\",\"office_id\":\"office2\"}]\n",
			},
		},
		{
			name: "customer get offices success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/offices", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				offices := []models.OfficeDTO{
					{BuildingID: "123", FloorNumber: 1, OfficeName: "kaushik office", OfficeID: "office1"},
				}
				mockOfficeService.EXPECT().ListOfficesByBuilding(gomock.Any(), "123").Return(offices, nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "[{\"building_id\":\"123\",\"floor_number\":1,\"office_name\":\"kaushik office\",\"office_id\":\"office1\"}]\n",
			},
		},
		{
			name: "get offices service error",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/123/offices", nil)
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().ListOfficesByBuilding(gomock.Any(), "123").Return(nil, errors.New("building not found"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to fetch offices\",\"code\":500}\n",
			},
		},
		{
			name: "empty offices list",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/api/v1/buildings/456/offices", nil)
					req.SetPathValue("buildingId", "456")
					return req
				}(),
			},
			mock: func() {
				offices := []models.OfficeDTO{}
				mockOfficeService.EXPECT().ListOfficesByBuilding(gomock.Any(), "456").Return(offices, nil)
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
			handler := &WebOfficeHandler{
				officeService: tt.fields.officeService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.GetOffices(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("GetOffices() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("GetOffices() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebOfficeHandler_AddOffice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeService := mocks.NewMockOfficeMgr(ctrl)

	type fields struct {
		officeService officeservice.OfficeMgr
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
			name: "add office success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/offices",
						bytes.NewBufferString(`{"office_name":"kaushik office","floor_number":1}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().AddOffice(gomock.Any(), "kaushik office", "123", 1).Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "",
			},
		},
		{
			name: "customer add office success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/offices",
						bytes.NewBufferString(`{"office_name":"kaushik office","floor_number":2}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().AddOffice(gomock.Any(), "kaushik office", "123", 2).Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "",
			},
		},
		{
			name: "invalid json request",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/offices",
						bytes.NewBufferString(`{"office_name":"kaushik office","floor_number":}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request payload\",\"code\":400}\n",
			},
		},
		{
			name: "add office service error",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/buildings/123/offices",
						bytes.NewBufferString(`{"office_name":"admin office","floor_number":3}`))
					req.SetPathValue("buildingId", "123")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().AddOffice(gomock.Any(), "admin office", "123", 3).Return(errors.New("office already exists"))
			},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"office already exists\",\"code\":400}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebOfficeHandler{
				officeService: tt.fields.officeService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.AddOffice(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("AddOffice() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("AddOffice() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestWebOfficeHandler_DeleteOffice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeService := mocks.NewMockOfficeMgr(ctrl)

	type fields struct {
		officeService officeservice.OfficeMgr
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
			name: "delete office success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/offices/office123", nil)
					req.SetPathValue("officeId", "office123")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().RemoveOffice(gomock.Any(), "office123").Return(nil)
			},
			want: want{
				status: http.StatusNoContent,
				body:   "",
			},
		},
		{
			name: "customer delete office success",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: customerCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/offices/office456", nil)
					req.SetPathValue("officeId", "office456")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().RemoveOffice(gomock.Any(), "office456").Return(nil)
			},
			want: want{
				status: http.StatusNoContent,
				body:   "",
			},
		},
		{
			name: "delete office service error",
			fields: fields{
				officeService: mockOfficeService,
			},
			args: args{
				ctx: adminCtx,
				w:   httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/api/v1/offices/office789", nil)
					req.SetPathValue("officeId", "office789")
					return req
				}(),
			},
			mock: func() {
				mockOfficeService.EXPECT().RemoveOffice(gomock.Any(), "office789").Return(errors.New("office not found"))
			},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"office not found\",\"code\":400}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := &WebOfficeHandler{
				officeService: tt.fields.officeService,
			}

			w := tt.args.w.(*httptest.ResponseRecorder)
			handler.DeleteOffice(tt.args.ctx, w, tt.args.r)

			if w.Code != tt.want.status {
				t.Errorf("DeleteOffice() status = %v, want %v", w.Code, tt.want.status)
			}

			if w.Body.String() != tt.want.body {
				t.Errorf("DeleteOffice() body = %v, want %v", w.Body.String(), tt.want.body)
			}
		})
	}
}
