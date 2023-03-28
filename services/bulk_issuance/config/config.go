package config

import (
	"errors"
	"log"

	"github.com/imroc/req"
	"github.com/jinzhu/configor"
)

func Initialize(fileName string) {
	err := configor.Load(&Config, fileName)
	if err != nil {
		panic("Unable to read configurations")
	}
	if Config.Keycloak.PublicKey == "" {
		updatePublicKeyFromKeycloak()
	}
}

var Config = struct {
	Registry struct {
		BaseUrl string `env:"REGISTRY_BASE_URL" yaml:"baseUrl" default:"http://localhost:8081/"`
	}
	Keycloak struct {
		PublicKey string `env:"KEYCLOAK_PUBLIC_KEY" yaml:"publicKey" default:""`
		Url       string `env:"KEYCLOAK_URL" yaml:"url" default:"http://keycloak:8080/auth"`
		Realm     string `env:"KEYCLOAK_REALM" yaml:"realm" default:"sunbird-rc"`
	}
}{}

func updatePublicKeyFromKeycloak() error {
	url := Config.Keycloak.Url + "/realms/" + Config.Keycloak.Realm
	log.Printf("Public key url : %v", url)
	resp, err := req.Get(url)
	if err != nil {
		return err
	}
	log.Printf("Got response %+v", resp.String())
	responseObject := map[string]interface{}{}
	if err := resp.ToJSON(&responseObject); err == nil {
		if publicKey, ok := responseObject["public_key"].(string); ok {
			Config.Keycloak.PublicKey = publicKey
		}
	}
	return errors.New("unable to get public key from keycloak")
}
