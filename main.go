package main

import (
	"ExchangeApp/config"
	"ExchangeApp/router"
)

func main() {
	config.InitConfig()

	r := router.SetRouter()

	port := config.AppConfig.App.Port

	if port == "" {
		port = ":8080"
	}
	r.Run(port)
	// listen and serve on 0.0.0.0:8080
}
