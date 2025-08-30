package authhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"go.uber.org/mock/gomock"
)

func TestNewWebAuthHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthenticationManager(ctrl)

	type args struct {
		service authservice.AuthenticationManager
	}
	tests := []struct {
		name string
		args args
		want *WebAuthHandler
	}{
		{
			name: "valid services",
			args: args{
				service: mockAuthService,
			},
			want: &WebAuthHandler{
				authServ: mockAuthService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebAuthHandler(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebAuthHandler_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthenticationManager(ctrl)

	type fields struct {
		authServ authservice.AuthenticationManager
	}
	type want struct {
		status int
		body   string
	}
	tests := []struct {
		name   string
		fields fields
		mock   func()
		req    any
		want   want
	}{
		{
			name: "valid login",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.LoginRequestDTO{
				Email:    "kaushik@a.com",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Login(gomock.Any()).Return("tokenxyz", nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "{\"jwt\":\"tokenxyz\"}",
			},
		},
		{
			name: "invalid json request body",
			fields: fields{
				authServ: mockAuthService,
			},
			req:  `{"email":"kaushik@a.com","password":}`,
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				//body:   "{\"message\":\"Invalid request body\",\"code\":400}\n",
			},
		},
		{
			name: "login service error",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.LoginRequestDTO{
				Email:    "kaushik@a.com",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Login(gomock.Any()).Return("", errors.New("invalid credentials"))
			},
			want: want{
				status: http.StatusInternalServerError,
				//body:   "{\"message\":\"Failed to login: invalid credentials\",\"code\":500}\n",
			},
		},
		{
			name: "admin login success",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.LoginRequestDTO{
				Email:    "admin@a.com",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Login(gomock.Any()).Return("admin_token_xyz", nil)
			},
			want: want{
				status: http.StatusOK,
				body:   "{\"jwt\":\"admin_token_xyz\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := WebAuthHandler{
				authServ: tt.fields.authServ,
			}

			var data []byte
			var err error

			if str, ok := tt.req.(string); ok {
				data = []byte(str)
			} else {
				data, err = json.Marshal(tt.req)
				if err != nil {
					t.Error(err)
				}
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(data))

			handler.Login(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.want.status {
				t.Errorf("Login() = %v, want %v", resp.StatusCode, tt.want.status)
			}

			//body, _ := io.ReadAll(resp.Body)
			//if string(body) != tt.want.body {
			//	t.Errorf("Login() = %v, want %v", string(body), tt.want.body)
			//}
		})
	}
}

func TestWebAuthHandler_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthenticationManager(ctrl)

	type fields struct {
		authServ authservice.AuthenticationManager
	}
	type want struct {
		status int
		body   string
	}
	tests := []struct {
		name   string
		fields fields
		mock   func()
		req    any
		want   want
	}{
		{
			name: "valid signup",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.RegisterRequestDTO{
				Name:     "kaushik",
				Email:    "kaushik@a.com",
				Office:   "office1",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "",
			},
		},
		{
			name: "invalid json request body",
			fields: fields{
				authServ: mockAuthService,
			},
			req:  `{"name":"kaushik","email":"kaushik@a.com","officeName":}`,
			mock: func() {},
			want: want{
				status: http.StatusBadRequest,
				body:   "{\"message\":\"Invalid request body\",\"code\":400}\n",
			},
		},
		{
			name: "signup service error",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.RegisterRequestDTO{
				Name:     "kaushik",
				Email:    "kaushik@a.com",
				Office:   "office1",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(errors.New("email already exists"))
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   "{\"message\":\"Failed to signup: email already exists\",\"code\":500}\n",
			},
		},
		{
			name: "admin signup",
			fields: fields{
				authServ: mockAuthService,
			},
			req: models.RegisterRequestDTO{
				Name:     "kaushik",
				Email:    "admin@a.com",
				Office:   "admin_office",
				Password: "123",
			},
			mock: func() {
				mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: want{
				status: http.StatusCreated,
				body:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			handler := WebAuthHandler{
				authServ: tt.fields.authServ,
			}

			var data []byte
			var err error

			if str, ok := tt.req.(string); ok {
				data = []byte(str)
			} else {
				data, err = json.Marshal(tt.req)
				if err != nil {
					t.Error(err)
				}
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(data))

			handler.Signup(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.want.status {
				t.Errorf("Signup() = %v, want %v", resp.StatusCode, tt.want.status)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != tt.want.body {
				t.Errorf("Signup() = %v, want %v", string(body), tt.want.body)
			}
		})
	}
}
