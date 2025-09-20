package floorhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebFloorHandler struct {
	floorService floorservice.FloorMgr
}

func NewWebFloorHandler(floorService floorservice.FloorMgr) *WebFloorHandler {
	return &WebFloorHandler{
		floorService: floorService,
	}
}

func (handler *WebFloorHandler) GetFloors(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	if userCtx.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	buildingId := r.PathValue("buildingId")

	floors, err := handler.floorService.GetFloorsByBuildingId(ctx, buildingId)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			// Message: "Failed to fetch floors",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	json.NewEncoder(w).Encode(floors)
}

// TODO: respect query params
func (handler *WebFloorHandler) AddFloor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	if userCtx.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	buildignId := r.PathValue("buildingId")

	var req struct {
		FloorNumber int `json:"floor_number"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.floorService.AddFloorByBuildingId(ctx, buildignId, req.FloorNumber)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			// Message: "Failed to add floor",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *WebFloorHandler) DeleteFloor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	if userCtx.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	buildignId := r.PathValue("buildingId")
	floorNumber, err := strconv.Atoi(r.PathValue("floorId"))
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid floor ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.floorService.DeleteFloor(ctx, buildignId, floorNumber)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			// Message: "Failed to delete floor",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
