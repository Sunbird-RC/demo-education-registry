package config

import "github.com/jinzhu/configor"

var Config = struct {
	Digilocker struct {
		TokenURL     string `env:"DIGILOCKER_SSO_TOKEN_URL" yaml:"digilocker_sso_token_url"`
		ClientId     string `env:"DIGILOCKER_SSO_CLIENT_ID" yaml:"digilocker_sso_client_id"`
		ClientSecret string `env:"DIGILOCKER_SSO_CLIENT_SECRET" yaml:"digilocker_sso_client_secret"`
		CodeVerifier string `env:"DIGILOCKER_SSO_CODE_VERIFIER" yaml:"digilocker_sso_code_verifier"`
		RedirectURL  string `env:"DIGILOCKER_SSO_REDIRECT_URL" yaml:"digilocker_sso_redirect_url"`
	}
	LogLevel string `env:"LOG_LEVEL" yaml:"log_level" default:"info"`
	Host     string `env:"HOST" yaml:"host" default:"0.0.0.0"`
	Port     string `env:"PORT" yaml:"port" default:"8085"`
	MODE     string `env:"MODE" yaml:"mode" default:"debug" `
}{}

func Init() {
	err := configor.Load(&Config, "./config/application-default.yml") //"config/application.yml"

	if err != nil {
		panic("Unable to read configurations")
	}

}
