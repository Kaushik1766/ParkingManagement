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
		"POST /auth/register":    app.AuthHandler.Signup,
		"POST /auth/login":       app.AuthHandler.Login,
		"GET /users":             authenticationmiddleware.AuthenticatedRoute(app.UserHandler.GetAllUsers),
		"PATCH /users/{userId}":  authenticationmiddleware.AuthenticatedRoute(app.UserHandler.UpdateProfile),
		"DELETE /users/{userId}": authenticationmiddleware.AuthenticatedRoute(app.UserHandler.DeleteProfile),
	}

	for route, handler := range routes {
		pathArr := strings.Split(route, " ")
		method := pathArr[0]
		path := basePath + pathArr[1]
		app.apiMux.HandleFunc(method+" "+path, handler)
	}
}
