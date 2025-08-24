package buildinghandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebBuildingHandler struct {
	buildingService buildingservice.BuildingMgr
}

func NewWebBuildingHandler(buildingService buildingservice.BuildingMgr) *WebBuildingHandler {
	return &WebBuildingHandler{
		buildingService: buildingService,
	}
}

func (handler *WebBuildingHandler) GetBuildings(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	if userCtx.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	query := r.URL.Query()

	var buildings []models.BuildingDTO
	var err error

	if query.Get("buildingId") != "" {
		building, err := handler.buildingService.GetBuildingByID(ctx, query.Get("buildingId"))
		if err != nil {
			customerrors.InternalServerError(w, customerrors.WebError{
				Message: "Failed to fetch building",
				Code:    http.StatusInternalServerError,
			})
			return
		}
		buildings = append(buildings, building)
	} else {
		buildings, err = handler.buildingService.GetAllBuildings(ctx)
	}

	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to fetch buildings",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	json.NewEncoder(w).Encode(buildings)
}

func (handler *WebBuildingHandler) AddBuilding(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	if userCtx.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req struct {
		Name string `json:"building_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request payload",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.buildingService.AddBuilding(ctx, req.Name)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to add building",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Building added successfully"})
}

func (handler *WebBuildingHandler) DeleteBuilding(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	err := handler.buildingService.DeleteBuildingByID(ctx, buildingId)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to delete building",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
