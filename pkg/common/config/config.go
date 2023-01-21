package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

type Config struct {
	BucketName          string
	StorageHost         string
	EmailFrom           string
	EmailPassword       string
	Host                string
	DBName              string
	MongoDBUri          string
	MqttUrl             *url.URL
	RedisAddr           string
	RedisPassword       string
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

		Cfg.BucketName = getenvStr("BUCKET_NAME")
		Cfg.StorageHost = getenvStr("STORAGE_HOST")

		Cfg.EmailFrom = getenvStr("EMAIL_FROM")
		Cfg.EmailPassword = getenvStr("EMAIL_PASSWORD")

		Cfg.DBName = getenvStr("DB_NAME")
		Cfg.MongoDBUri = getenvStr("MONGODB_URI")

		var err error
		Cfg.MqttUrl, err = url.Parse(getenvStr("MQTT_URL"))
		if err != nil {
			log.Fatal(err)
		}

		Cfg.RedisAddr = getenvStr("REDIS_ADDR")
		Cfg.RedisPassword = getenvStr("REDIS_PASSWORD")

		Cfg.JWTAccessSecretKey = getenvStr("JWT_ACCESS_SECRET_KEY")
		Cfg.JWTRefreshSecretKey = getenvStr("JWT_REFRESH_SECRET_KEY")
	}
}
