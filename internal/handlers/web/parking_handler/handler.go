package parkinghandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type WebParkingHandler struct {
	parkingService parkinghistoryservice.ParkingHistoryMgr
	vehicleService vehicleservice.VehicleMgr
}

func NewWebParkingHandler(parkingService parkinghistoryservice.ParkingHistoryMgr, vehicleService vehicleservice.VehicleMgr) *WebParkingHandler {
	return &WebParkingHandler{
		parkingService: parkingService,
		vehicleService: vehicleService,
	}
}

func (handler *WebParkingHandler) GetParkings(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	var startTime, endTime time.Time
	var err error

	startTimeStr := query.Get("startTime")
	if startTimeStr == "" {
		startTime = time.Now().AddDate(0, -1, 0)
	} else {
		startTime, err = time.Parse(time.RFC3339, query.Get("startTime"))
		if err != nil {
			customerrors.BadRequestError(w, customerrors.WebError{
				Message: "Invalid startTime format. Use RFC3339 format.",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	endTimeStr := query.Get("endTime")
	if endTimeStr == "" {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(time.RFC3339, query.Get("endTime"))
		if err != nil {
			customerrors.BadRequestError(w, customerrors.WebError{
				Message: "Invalid endTime format. Use RFC3339 format.",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}
	// TODO: add option for admin to get parking history for all users
	parkings, err := handler.parkingService.GetParkingHistory(ctx, startTime, endTime)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to fetch parking history: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parkings)
}

func (handler *WebParkingHandler) AddParking(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var parkingRequest struct {
		NumberPlate string `json:"numberplate"`
	}

	err := json.NewDecoder(r.Body).Decode(&parkingRequest)
	if err != nil {
		customerrors.BadRequestError(w, customerrors.WebError{
			Message: "Invalid request payload: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	ticketId, err := handler.vehicleService.Park(ctx, parkingRequest.NumberPlate)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to park vehicle: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		TicketID string `json:"ticketId"`
	}{
		TicketID: ticketId,
	})
}

func (handler *WebParkingHandler) UnparkVehicle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	numberPlate := r.PathValue("numberplate")

	err := handler.vehicleService.UnparkByNumberPlate(ctx, numberPlate)
	if err != nil {
		customerrors.InternalServerError(w, customerrors.WebError{
			Message: "Failed to unpark vehicle: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
