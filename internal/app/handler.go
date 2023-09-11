package app

import (
	"github.com/Tap-Team/kurilka/internal/middleware"
	"github.com/gorilla/mux"
)

func Router(
	vk_secret_key string,
) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)
	r.Use(middleware.VK)
	r.Use(middleware.LaunchParams(vk_secret_key))
	return r
}
