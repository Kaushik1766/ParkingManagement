package userhandler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebUserHandler struct {
	userService userservice.UserManager
}

func NewWebUserHandler(service userservice.UserManager) *WebUserHandler {
	return &WebUserHandler{
		userService: service,
	}
}

func (handler *WebUserHandler) GetAllUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctxUser := ctx.Value(constants.User).(models.UserJwt)

	queryParams := r.URL.Query()

	var err error
	var res []models.UserDTO
	if ctxUser.Role == roles.Admin {
		if queryParams.Get("userId") == "" {
			res, err = handler.userService.GetAllUsers(ctx)
			if err != nil {
				customerrors.InternalServerError(w, customerrors.WebError{
					Message: "Failed to fetch users: " + err.Error(),
					Code:    http.StatusInternalServerError,
				})
				return
			}
		} else {
			u, err := handler.userService.GetUserById(ctx, queryParams.Get("userId"))
			if err != nil {
				customerrors.InternalServerError(w, customerrors.WebError{
					Message: "Failed to fetch user: " + err.Error(),
					Code:    http.StatusInternalServerError,
				})
				return
			}
			res = append(res, u)
		}
	} else {
		u, err := handler.userService.GetUserById(ctx, ctxUser.Subject)
		if err != nil {
			customerrors.InternalServerError(w, customerrors.WebError{
				Message: "Failed to fetch user: " + err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}
		res = append(res, u)
	}
	json.NewEncoder(w).Encode(res)
}

func (handler *WebUserHandler) UpdateProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.PathValue("userId")
	data, _ := io.ReadAll(r.Body)
	//if err != nil {
	//	customerrors.InternalServerError(w, customerrors.WebError{
	//		Message: "Failed to read request body",
	//		Code:    http.StatusInternalServerError,
	//	})
	//	return
	//}

	var req models.UpdateUserDTO

	err := json.Unmarshal(data, &req)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.userService.UpdateProfile(ctx, userId, req)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to update profile: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *WebUserHandler) DeleteProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.PathValue("userId")

	err := handler.userService.DeleteProfile(ctx, userId)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to delete profile: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
