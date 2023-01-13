package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	cMqtt "github.com/jovanfrandika/smartbox-backend/pkg/common/mqtt"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	deviceService "github.com/jovanfrandika/smartbox-backend/pkg/device/service"
	rFriendshipMongo "github.com/jovanfrandika/smartbox-backend/pkg/friendship/repository/mongo"
	friendshipService "github.com/jovanfrandika/smartbox-backend/pkg/friendship/service"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	parcelService "github.com/jovanfrandika/smartbox-backend/pkg/parcel/service"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
	userService "github.com/jovanfrandika/smartbox-backend/pkg/user/service"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
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

	storageClient, err := storage.NewClient(context.TODO(), option.WithCredentialsFile("files/service-account.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer storageClient.Close()

	userDb := rUserMongo.New(db)
	friendshipDb := rFriendshipMongo.New(db)
	parcelDb := rParcelMongo.New(db)
	deviceDb := rDeviceMongo.New(db)

	mqttClient := cMqtt.Init("stancyzk", *config.Cfg)
	deviceMq := rDeviceMqtt.New(&mqttClient)

	r := chi.NewRouter()

	userRouter := chi.NewRouter()
	userService.Init(&userDb, userRouter, config.Cfg)

	deviceRouter := chi.NewRouter()
	deviceService.Init(&deviceDb, deviceRouter, config.Cfg)

	parcelRouter := chi.NewRouter()
	parcelService.Init(parcelService.InitInput{
		ParcelDb:      &parcelDb,
		UserDb:        &userDb,
		DeviceDb:      &deviceDb,
		DeviceMq:      &deviceMq,
		Config:        config.Cfg,
		Router:        parcelRouter,
		StorageClient: storageClient,
	})

	friendshipRouter := chi.NewRouter()
	friendshipService.Init(&friendshipDb, &userDb, friendshipRouter, config.Cfg)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	r.Mount("/user", userRouter)
	r.Mount("/friendship", friendshipRouter)
	r.Mount("/device", deviceRouter)
	r.Mount("/parcel", parcelRouter)

	http.ListenAndServe(":8000", r)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
