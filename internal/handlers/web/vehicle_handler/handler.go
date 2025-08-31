package vehiclehandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebVehicleHandler struct {
	vehicleService vehicleservice.VehicleMgr
	userService    userservice.UserManager
}

func NewWebVehicleHandler(vehicleService vehicleservice.VehicleMgr, userService userservice.UserManager) *WebVehicleHandler {
	return &WebVehicleHandler{
		vehicleService: vehicleService,
		userService:    userService,
	}
}

func (handler *WebVehicleHandler) GetVehicles(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctxUser := ctx.Value(constants.User).(models.UserJwt)

	var vehicles []models.VehicleDTO
	var err error

	if ctxUser.Role == roles.Admin {
		queryParams := r.URL.Query()
		if queryParams.Get("userId") != "" {
			vehicles, err = handler.userService.GetVehiclesByUserId(ctx, queryParams.Get("userId"))
			if err != nil {
				customerrors.SendError(w, err)
				return
			}
		}
	} else {
		vehicles, err = handler.userService.GetRegisteredVehicles(ctx)
		if err != nil {
			customerrors.InternalServerError(w, customerrors.WebError{
				Message: "Failed to fetch vehicles",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

func (handler *WebVehicleHandler) RegisterVehicle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// ctxUser := ctx.Value(constants.User).(models.UserJwt)

	var vehicleReq models.AddVehicleDTO
	err := json.NewDecoder(r.Body).Decode(&vehicleReq)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request payload",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if vehicleReq.VehicleType != 0 && vehicleReq.VehicleType != 1 {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid vehicle type",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.userService.RegisterVehicle(ctx, vehicleReq.NumberPlate, vehicletypes.VehicleType(vehicleReq.VehicleType))
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *WebVehicleHandler) RemoveVehicle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// ctxUser := ctx.Value(constants.User).(models.UserJwt)

	numberplate := r.PathValue("numberplate")

	err := handler.userService.UnregisterVehicle(ctx, numberplate)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
