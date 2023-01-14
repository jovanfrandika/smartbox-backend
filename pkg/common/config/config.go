package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

type Config struct {
	MqttUrl             *url.URL
	BucketName          string
	DBName              string
	MongoDBUri          string
	JWTAccessSecretKey  string
	JWTRefreshSecretKey string
}

var Cfg *Config

func getenvStr(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal(fmt.Sprintf("%s is empty", key))
	}
	return v
}

func Init() {
	if Cfg == nil {
		Cfg = &Config{}

		var err error
		Cfg.MqttUrl, err = url.Parse(getenvStr("MQTT_URL"))
		if err != nil {
			log.Fatal(err)
		}
		Cfg.BucketName = getenvStr("BUCKET_NAME")
		Cfg.DBName = getenvStr("DB_NAME")
		Cfg.MongoDBUri = getenvStr("MONGODB_URI")
		Cfg.JWTAccessSecretKey = getenvStr("JWT_ACCESS_SECRET_KEY")
		Cfg.JWTRefreshSecretKey = getenvStr("JWT_REFRESH_SECRET_KEY")
	}
}