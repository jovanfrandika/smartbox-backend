package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	cMqtt "github.com/jovanfrandika/smartbox-backend/pkg/common/mqtt"
	deviceConsumer "github.com/jovanfrandika/smartbox-backend/pkg/device/consumer"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	parcelTravelConsumer "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/consumer"
	rParcelTravelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/repository/mongo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	config.Init()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.Cfg.MongoDBUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(config.Cfg.DBName)

	parcelDb := rParcelMongo.New(db)
	parcelTravelDb := rParcelTravelMongo.New(db)
	deviceDb := rDeviceMongo.New(db)

	mqttClient := cMqtt.Init("jester", *config.Cfg)
	deviceConsumer.Init(&deviceDb, &mqttClient, config.Cfg)
	parcelTravelConsumer.Init(&parcelTravelDb, &parcelDb, &deviceDb, &mqttClient, config.Cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
