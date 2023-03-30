package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sso-service/config"
	"sso-service/services"
)

type DigilockerOIDC struct {
}

func (d DigilockerOIDC) GetToken(c *gin.Context) {
	log.Info("Get token API triggered")
	code := c.PostForm("code")
	url := config.Config.Digilocker.TokenURL
	clientId := config.Config.Digilocker.ClientId
	clientSecret := config.Config.Digilocker.ClientSecret
	codeVerifier := config.Config.Digilocker.CodeVerifier
	redirectURL := config.Config.Digilocker.RedirectURL
	payload := fmt.Sprintf("code=%s&grant_type=authorization_code&client_id=%s&client_secret=%s&code_verifier=%s"+
		"&redirect_uri=%s", code, clientId, clientSecret, codeVerifier, redirectURL)
	log.Debugf("Digilocker API token payload {%s}", payload)
	ssoResponseData := services.CallPOST(url, payload)
	if ssoResponseData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed fetching digilocker token",
		})
	} else {
		response := gin.H{
			"access_token": ssoResponseData["access_token"].(string) + ":" + ssoResponseData["id_token"].(string),
			"id_token":     ssoResponseData["id_token"],
			"expires_in":   ssoResponseData["expires_in"],
			"token_type":   ssoResponseData["token_type"],
			"scope":        ssoResponseData["scope"],
		}
		log.Debugf("Token response, %v", response)
		log.Info("Get token API completed")
		c.JSON(http.StatusOK, response)
	}
}
