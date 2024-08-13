package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type Services struct {
	Client services.Client
	Auth   services.Auth
}

type Routes struct {
	services Services
	repo     *repo.CRUDRepo
}

func NewRoutes(services Services, repo *repo.CRUDRepo) *Routes {
	return &Routes{
		services: services,
		repo:     repo,
	}
}

func (r *Routes) NewRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	NewAuthRoutes(api, r.services.Auth)

	ws := router.PathPrefix("/ws").Subrouter()
	NewClientRoutes(ws, r.services.Client)

	// api.Handle("/test", middleware.UserProtectionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("hello")
	// }), r.repo))
	return router
}
