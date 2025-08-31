package userhandler

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
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestNewWebUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserManager(ctrl)

	tests := []struct {
		name    string
		service *mocks.MockUserManager
		want    *WebUserHandler
	}{
		{
			name:    "successful creation",
			service: mockUserService,
			want: &WebUserHandler{
				userService: mockUserService,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWebUserHandler(tt.service)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebUserHandler_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockUserManager)
		setupContext   func() context.Context
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "admin gets all users without userid query",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetAllUsers(gomock.Any()).Return([]models.UserDTO{
					{
						UserId: "user-123",
						Name:   "kaushik",
						Email:  "kaushik@a.com",
						Role:   "Customer",
						Office: "office-1",
					},
					{
						UserId: "admin-456",
						Name:   "admin",
						Email:  "admin@a.com",
						Role:   "Admin",
						Office: "office-1",
					},
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "admin-456"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/users", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "admin gets specific user by userid query",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetUserById(gomock.Any(), "user-123").Return(models.UserDTO{
					UserId: "user-123",
					Name:   "kaushik",
					Email:  "kaushik@a.com",
					Role:   "Customer",
					Office: "office-1",
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "admin-456"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/users?userId=user-123", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "customer gets own profile",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetUserById(gomock.Any(), "user-123").Return(models.UserDTO{
					UserId: "user-123",
					Name:   "kaushik",
					Email:  "kaushik@a.com",
					Role:   "Customer",
					Office: "office-1",
				}, nil)
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "user-123"},
					Email:            "kaushik@a.com",
					Role:             roles.Customer,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/users", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "admin get all users service error",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetAllUsers(gomock.Any()).Return(nil, errors.New("database error"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "admin-456"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/users", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to fetch users: database error","code":500}`,
		},
		{
			name: "admin get user by id service error",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetUserById(gomock.Any(), "user-123").Return(models.UserDTO{}, errors.New("user not found"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "admin-456"},
					Email:            "admin@a.com",
					Role:             roles.Admin,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/users?userId=user-123", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to fetch user: user not found","code":500}`,
		},
		{
			name: "customer get own profile service error",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().GetUserById(gomock.Any(), "user-123").Return(models.UserDTO{}, errors.New("user not found"))
			},
			setupContext: func() context.Context {
				userJwt := models.UserJwt{
					RegisteredClaims: jwt.RegisteredClaims{ID: "user-123"},
					Email:            "kaushik@a.com",
					Role:             roles.Customer,
					Office:           "office-1",
				}
				return context.WithValue(context.Background(), constants.User, userJwt)
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/users", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to fetch user: user not found","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockUserManager(ctrl)
			tt.setupMock(mockService)

			handler := NewWebUserHandler(mockService)
			w := httptest.NewRecorder()
			ctx := tt.setupContext()
			req := tt.setupRequest()

			handler.GetAllUsers(ctx, w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetAllUsers() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("GetAllUsers() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("GetAllUsers() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}

			// For successful cases, verify the response structure
			if tt.expectedStatus == http.StatusOK && tt.expectedBody == "" {
				var users []models.UserDTO
				if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
					t.Errorf("GetAllUsers() failed to unmarshal response: %v", err)
				}
				if len(users) == 0 {
					t.Errorf("GetAllUsers() returned empty users array")
				}
			}
		})
	}
}

func TestWebUserHandler_UpdateProfile(t *testing.T) {
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
			name: "successful profile update",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().UpdateProfile(gomock.Any(), "user-123", models.UpdateUserDTO{
					Name:   "kaushik updated",
					Email:  "kaushik_new@a.com",
					Office: "office-2",
				}).Return(nil)
			},
			setupRequest: func() *http.Request {
				updateReq := models.UpdateUserDTO{
					Name:   "kaushik updated",
					Email:  "kaushik_new@a.com",
					Office: "office-2",
				}
				body, _ := json.Marshal(updateReq)
				req := httptest.NewRequest("PUT", "/users/user-123", bytes.NewBuffer(body))
				req.SetPathValue("userId", "user-123")
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid json request body",
			setupMock: func(mockService *mocks.MockUserManager) {
				// No expectations as the request should fail before reaching the service
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("PUT", "/users/user-123", bytes.NewBuffer([]byte("invalid json")))
				req.SetPathValue("userId", "user-123")
				return req
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid request body","code":400}`,
		},
		{
			name: "service error during update",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().UpdateProfile(gomock.Any(), "user-123", models.UpdateUserDTO{
					Name: "kaushik",
				}).Return(errors.New("database error"))
			},
			setupRequest: func() *http.Request {
				updateReq := models.UpdateUserDTO{
					Name: "kaushik",
				}
				body, _ := json.Marshal(updateReq)
				req := httptest.NewRequest("PUT", "/users/user-123", bytes.NewBuffer(body))
				req.SetPathValue("userId", "user-123")
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to update profile: database error","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockUserManager(ctrl)
			tt.setupMock(mockService)

			handler := NewWebUserHandler(mockService)
			w := httptest.NewRecorder()
			req := tt.setupRequest()

			handler.UpdateProfile(context.Background(), w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateProfile() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("UpdateProfile() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("UpdateProfile() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}
		})
	}
}

func TestWebUserHandler_DeleteProfile(t *testing.T) {
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
			name: "successful profile deletion",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().DeleteProfile(gomock.Any(), "user-123").Return(nil)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/users/user-123", nil)
				req.SetPathValue("userId", "user-123")
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error during deletion",
			setupMock: func(mockService *mocks.MockUserManager) {
				mockService.EXPECT().DeleteProfile(gomock.Any(), "user-123").Return(errors.New("database error"))
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("DELETE", "/users/user-123", nil)
				req.SetPathValue("userId", "user-123")
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Failed to delete profile: database error","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMockUserManager(ctrl)
			tt.setupMock(mockService)

			handler := NewWebUserHandler(mockService)
			w := httptest.NewRecorder()
			req := tt.setupRequest()

			handler.DeleteProfile(context.Background(), w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("DeleteProfile() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if w.Body.String() != tt.expectedBody+"\n" {
					t.Errorf("DeleteProfile() body = %v, want %v", w.Body.String(), tt.expectedBody+"\n")
				}
			}

			// Verify content type is set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("DeleteProfile() Content-Type = %v, want %v", w.Header().Get("Content-Type"), "application/json")
			}
		})
	}
}
