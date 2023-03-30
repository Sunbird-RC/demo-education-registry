package server

import (
	"github.com/gin-gonic/gin"
	"sso-service/config"
	"sso-service/controllers"
	"sso-service/middlewares"
)

func NewRouter() *gin.Engine {
	gin.SetMode(config.Config.MODE)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{
		digilockerSSOGroup := v1.Group("oidc/digilocker")
		{
			var oidc controllers.OIDC = controllers.DigilockerOIDC{}
			digilockerSSOGroup.POST("/token", oidc.GetToken)
		}
	}
	return router

}
