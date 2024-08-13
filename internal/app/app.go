package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/middleware"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/internal/routes"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

func Run(config config.Config) {
	crudRepo := repo.NewCRUDRepo(config)

	services := routes.Services{
		Auth:   services.NewAuthService(*crudRepo),
		Client: services.NewClientService(*crudRepo),
	}

	routes := routes.NewRoutes(services, crudRepo)
	router := routes.NewRouter()

	handler := middleware.CorsMiddleware(router)

	server := &http.Server{
		Addr:         config.Site.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(config.Site.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Site.WriteTimeout) * time.Second,
	}

	fmt.Println("Listen and serve", config.Site.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
