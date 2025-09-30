package officehandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebOfficeHandler struct {
	officeService officeservice.OfficeMgr
}

func NewWebOfficeHandler(officeService officeservice.OfficeMgr) *WebOfficeHandler {
	return &WebOfficeHandler{
		officeService: officeService,
	}
}

func (handler *WebOfficeHandler) GetAllOffices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	offices, err := handler.officeService.GetAllOfficeNames(context.Background())
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to fetch offices",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	json.NewEncoder(w).Encode(offices)
}

func (handler *WebOfficeHandler) GetOffices(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	buildingId := r.PathValue("buildingId")

	offices, err := handler.officeService.ListOfficesByBuilding(ctx, buildingId)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to fetch offices",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	json.NewEncoder(w).Encode(offices)
}

func (handler *WebOfficeHandler) AddOffice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	buildingId := r.PathValue("buildingId")
	var officeReq models.OfficeDTO
	err := json.NewDecoder(r.Body).Decode(&officeReq)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request payload",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = handler.officeService.AddOffice(ctx, officeReq.OfficeName, buildingId, officeReq.FloorNumber)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *WebOfficeHandler) DeleteOffice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	officeId := r.PathValue("officeId")

	err := handler.officeService.RemoveOffice(ctx, officeId)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
