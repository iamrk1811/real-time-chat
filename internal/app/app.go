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

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Error upgrading:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Error reading:", err)
// 			return
// 		}
// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			fmt.Println("Error writing:", err)
// 			return
// 		}
// 	}
// }

func Run(config config.Config) {
	crudRepo := repo.NewCRUDRepo(config)

	services := routes.Services{
		User: services.NewUserService(*crudRepo),
	}

	routes := routes.NewRoutes(services)
	router := routes.NewRouter()

	server := &http.Server{
		Addr:         config.Site.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(config.Site.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Site.WriteTimeout) * time.Second,
	}

	fmt.Println("Listen and serve", config.Site.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
