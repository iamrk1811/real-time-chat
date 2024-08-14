package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/middleware"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type Services struct {
	Client services.Client
	Auth   services.Auth
}

type Routes struct {
	Services Services
	Repo     *repo.CRUDRepo
}

func NewRoutes(services Services, repo *repo.CRUDRepo) *Routes {
	return &Routes{
		Services: services,
		Repo:     repo,
	}
}

func (r *Routes) NewRouter(config *config.Config) http.Handler {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	NewAuthRoutes(api, r.Services.Auth)
	NewClientRoutes(api, r.Services.Client, config)

	ws := router.PathPrefix("/ws").Subrouter()
	_ = ws

	// api.Handle("/test", middleware.UserProtectionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("hello")
	// }), r.repo))

	handler := middleware.CorsMiddleware(router)
	handler = middleware.SessionProtection(handler, r.Repo, config)
	return handler
}
