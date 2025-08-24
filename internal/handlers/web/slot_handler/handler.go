package slothandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebSlotHandler struct {
	slotService slotservice.SlotMgr
}

func NewWebSlotHandler(slotService slotservice.SlotMgr) *WebSlotHandler {
	return &WebSlotHandler{
		slotService: slotService,
	}
}

func (handler *WebSlotHandler) GetSlots(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctxUser := ctx.Value(constants.User).(models.UserJwt)

	if ctxUser.Role != roles.Admin {
		customerrors.UnauthorizedError(w, customerrors.WebError{
			Message: "Unauthorized access",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	buildingId := r.PathValue("buildingId")
	floorNumber, err := strconv.Atoi(r.PathValue("floorId"))
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid floor ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	slots, err := handler.slotService.GetSlotsByFloor(ctx, buildingId, floorNumber)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to fetch slots",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	json.NewEncoder(w).Encode(slots)
}
