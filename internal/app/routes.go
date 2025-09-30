package app

import (
	"fmt"
	"net/http"
	"strings"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	loggingmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/logging_middleware"
)

var routes map[string]func(w http.ResponseWriter, r *http.Request)

var basePath = "/api/v1"

func (app *App) registerRoutes() {
	authMiddleware := authenticationmiddleware.AuthenticatedRoute
	loggingMiddleware := loggingmiddleware.LoggingMiddleware

	routes = map[string]func(w http.ResponseWriter, r *http.Request){
		"POST /auth/register":                                app.AuthHandler.Signup,
		"POST /auth/login":                                   app.AuthHandler.Login,
		"GET /users":                                         authMiddleware(app.UserHandler.GetAllUsers),
		"PATCH /users/{userId}":                              authMiddleware(app.UserHandler.UpdateProfile),
		"DELETE /users/{userId}":                             authMiddleware(app.UserHandler.DeleteProfile),
		"GET /buildings":                                     authMiddleware(app.BuildingHandler.GetBuildings),
		"POST /buildings":                                    authMiddleware(app.BuildingHandler.AddBuilding),
		"DELETE /buildings/{buildingId}":                     authMiddleware(app.BuildingHandler.DeleteBuilding),
		"GET /buildings/{buildingId}/floors":                 authMiddleware(app.FloorHandler.GetFloors),
		"POST /buildings/{buildingId}/floors":                authMiddleware(app.FloorHandler.AddFloor),
		"GET /buildings/{buildingId}/offices":                authMiddleware(app.OfficeHandler.GetOffices),
		"GET /offices":                                       app.OfficeHandler.GetAllOffices,
		"POST /buildings/{buildingId}/offices":               authMiddleware(app.OfficeHandler.AddOffice),
		"DELETE /buildings/{buildingId}/offices/{officeId}":  authMiddleware(app.OfficeHandler.DeleteOffice),
		"DELETE /buildings/{buildingId}/floors/{floorId}":    authMiddleware(app.FloorHandler.DeleteFloor),
		"GET /buildings/{buildingId}/floors/{floorId}/slots": authMiddleware(app.SlotHandler.GetSlots),
		"GET /vehicles":                                      authMiddleware(app.VehicleHandler.GetVehicles),
		"POST /vehicles":                                     authMiddleware(app.VehicleHandler.RegisterVehicle),
		"DELETE /vehicles/{vehicleId}":                       authMiddleware(app.VehicleHandler.RemoveVehicle),
		"POST /parkings":                                     authMiddleware(app.ParkingHandler.AddParking),
		"GET /parkings":                                      authMiddleware(app.ParkingHandler.GetParkings),
		"PATCH /parkings/{numberplate}/unpark":               authMiddleware(app.ParkingHandler.UnparkVehicle),
	}

	for route, handler := range routes {
		pathArr := strings.Split(route, " ")
		method := pathArr[0]
		path := basePath + pathArr[1]
		app.apiMux.HandleFunc(method+" "+path, loggingMiddleware(handler))
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("enableCORS")

		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allow common HTTP methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allow common headers

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
