package app

import (
	"github.com/Tap-Team/kurilka/internal/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logger)
	return r
}
