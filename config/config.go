package config

import (
	"os"
)

type AppConfig struct {
	Dbuser          string
	Dbpass          string
	Dbhost          string
	Dbport          string
	Dbname          string
	LibDir          string
	TempStoragePath string
	IPPK            string
	KeycloakUrl     string
	ClientName      string
	Realm           string
	AdminUsername   string
	AdminPassword   string
}

func GetEnv() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.Dbuser = os.Getenv("DB_USER")
	appConfig.Dbpass = os.Getenv("DB_PASS")
	appConfig.Dbhost = os.Getenv("DB_HOST")
	appConfig.Dbport = os.Getenv("DB_PORT")
	appConfig.Dbname = os.Getenv("DB_NAME")
	appConfig.LibDir = os.Getenv("LIB_DIR")
	appConfig.IPPK = os.Getenv("IPPK")
	appConfig.ClientName = os.Getenv("CLIENT_NAME")
	appConfig.Realm = os.Getenv("REALM")
	appConfig.KeycloakUrl = os.Getenv("KEYCLOAK_URL")
	appConfig.AdminUsername = os.Getenv("ADMIN_USERNAME")
	appConfig.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	return appConfig
}
