package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type Services struct {
	User services.User
}

type Routes struct {
	services Services
}

func NewRoutes(services Services) *Routes{
	return &Routes{
		services: services,
	}
}


func (r *Routes) NewRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()

	NewUserRoutes(api, r.services.User)
	NewChatRoutes(api)

	return router
}