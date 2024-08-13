package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/internal/routes"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override, Content-Type, Accept		, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If this is a preflight request, then stop here
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Run(config config.Config) {
	crudRepo := repo.NewCRUDRepo(config)

	services := routes.Services{
		Auth:   services.NewAuthService(*crudRepo),
		Client: services.NewClientService(*crudRepo),
	}

	routes := routes.NewRoutes(services)
	router := routes.NewRouter()

	handler := corsMiddleware(router)

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
