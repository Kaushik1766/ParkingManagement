package app

import (
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	db  *gorm.DB
	mux *http.ServeMux
}
