package app

import (
	"net/http"
	"strings"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
)

var routes map[string]func(w http.ResponseWriter, r *http.Request)

var basePath = "/api/v1"

func (app *App) registerRoutes() {
	routes = map[string]func(w http.ResponseWriter, r *http.Request){
		"POST /auth/register":                                app.AuthHandler.Signup,
		"POST /auth/login":                                   app.AuthHandler.Login,
		"GET /users":                                         authenticationmiddleware.AuthenticatedRoute(app.UserHandler.GetAllUsers),
		"PATCH /users/{userId}":                              authenticationmiddleware.AuthenticatedRoute(app.UserHandler.UpdateProfile),
		"DELETE /users/{userId}":                             authenticationmiddleware.AuthenticatedRoute(app.UserHandler.DeleteProfile),
		"GET /buildings":                                     authenticationmiddleware.AuthenticatedRoute(app.BuildingHandler.GetBuildings),
		"POST /buildings":                                    authenticationmiddleware.AuthenticatedRoute(app.BuildingHandler.AddBuilding),
		"DELETE /buildings/{buildingId}":                     authenticationmiddleware.AuthenticatedRoute(app.BuildingHandler.DeleteBuilding),
		"GET /buildings/{buildingId}/floors":                 authenticationmiddleware.AuthenticatedRoute(app.FloorHandler.GetFloors),
		"POST /buildings/{buildingId}/floors":                authenticationmiddleware.AuthenticatedRoute(app.FloorHandler.AddFloor),
		"GET /buildings/{buildingId}/offices":                authenticationmiddleware.AuthenticatedRoute(app.OfficeHandler.GetOffices),
		"POST /buildings/{buildingId}/offices":               authenticationmiddleware.AuthenticatedRoute(app.OfficeHandler.AddOffice),
		"DELETE /buildings/{buildingId}/offices/{officeId}":  authenticationmiddleware.AuthenticatedRoute(app.OfficeHandler.DeleteOffice),
		"DELETE /buildings/{buildingId}/floors/{floorId}":    authenticationmiddleware.AuthenticatedRoute(app.FloorHandler.DeleteFloor),
		"GET /buildings/{buildingId}/floors/{floorId}/slots": authenticationmiddleware.AuthenticatedRoute(app.SlotHandler.GetSlots),
		"GET /vehicles":                                      authenticationmiddleware.AuthenticatedRoute(app.VehicleHandler.GetVehicles),
		"POST /vehicles":                                     authenticationmiddleware.AuthenticatedRoute(app.VehicleHandler.RegisterVehicle),
	}

	for route, handler := range routes {
		pathArr := strings.Split(route, " ")
		method := pathArr[0]
		path := basePath + pathArr[1]
		app.apiMux.HandleFunc(method+" "+path, handler)
	}
}
