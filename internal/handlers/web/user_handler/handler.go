package userhandler

import userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"

type WebUserHandler struct {
	userService userservice.UserManager
}
