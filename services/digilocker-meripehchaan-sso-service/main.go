package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sso-service/config"
	"sso-service/server"
)

func main() {
	config.Init()
	ll, err := log.ParseLevel(config.Config.LogLevel)
	if err != nil {
		fmt.Print("Failed parsing log level")
	}
	log.SetLevel(ll)
	log.Info("Log level: %s", config.Config.LogLevel)
	server.Init()
}
