package main

import (
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/app"
)

func main() {
	config := config.NewConfig()
	app.Run(*config)
}
