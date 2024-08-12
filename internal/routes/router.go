package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type Services struct {
	Client services.Client
	Auth   services.Auth
}

type Routes struct {
	services Services
}

func NewRoutes(services Services) *Routes {
	return &Routes{
		services: services,
	}
}

func (r *Routes) NewRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	NewAuthRoutes(api, r.services.Auth)

	ws := router.PathPrefix("/ws").Subrouter()
	NewClientRoutes(ws, r.services.Client)

	return router
}
