package settings

import (
	"os"

	"github.com/ole-larsen/uploader/service/log"
	"github.com/spf13/viper"
)

var Settings = initSettings()

type settings struct {
	XToken  string
	APP     string
	Port    string
	Secret  string
	PG      PG
	PGSQL   string
	UseHash bool
	UseDB   bool
}

type PG struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func initSettings() settings {
	var ss settings
	logger := log.NewLogger()

	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		logger.Errorln("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	app, ok := viper.Get("APP_NAME").(string)
	if !ok {
		app = os.Getenv("APP_NAME")
	}
	ss.APP = app

	secret, ok := viper.Get("SESSION_SECRET").(string)
	if !ok {
		secret = os.Getenv("SESSION_SECRET")
	}
	ss.Secret = secret

	xToken, ok := viper.Get("X_TOKEN").(string)
	if !ok {
		xToken = os.Getenv("X_TOKEN")
	}
	ss.XToken = xToken

	port, ok := viper.Get("PORT").(string)
	if !ok {
		port = os.Getenv("PORT")
	}
	ss.Port = port

	dbHost := viper.GetString("DB_SQL_HOST")
	dbPort := viper.GetString("DB_SQL_PORT")
	dbUsername := viper.GetString("DB_SQL_USERNAME")
	dbPassword := viper.GetString("DB_SQL_PASSWORD")
	db := viper.GetString("DB_SQL_DATABASE")

	if dbHost == "" {
		dbHost = os.Getenv("DB_SQL_HOST")
	}
	if dbPort == "" {
		dbPort = os.Getenv("DB_SQL_PORT")
	}
	if dbUsername == "" {
		dbUsername = os.Getenv("DB_SQL_USERNAME")
	}
	if dbPassword == "" {
		dbPassword = os.Getenv("DB_SQL_PASSWORD")
	}
	if db == "" {
		db = os.Getenv("DB_SQL_DATABASE")
	}
	ss.PG = PG{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUsername,
		Password: dbPassword,
		Database: db,
	}
	pgsql := "postgres://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + db + "?sslmode=disable"
	ss.PGSQL = pgsql

	ss.UseHash = viper.GetBool("USE_HASH")
	ss.UseDB = viper.GetBool("USE_DB")

	logger.Println("load settings done √")

	return ss
}
