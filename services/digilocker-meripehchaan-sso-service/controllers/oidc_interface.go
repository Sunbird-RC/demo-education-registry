package controllers

import "github.com/gin-gonic/gin"

type OIDC interface {
	GetToken(c *gin.Context)
}
