package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sso-service/config"
)

func Init() {
	r := NewRouter()
	err := r.Run(fmt.Sprintf("%s:%s", config.Config.Host, config.Config.Port))
	if err != nil {
		log.Error("Failed starting the router")
	}
}
