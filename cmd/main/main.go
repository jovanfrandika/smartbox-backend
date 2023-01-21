package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v9"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	cEmail "github.com/jovanfrandika/smartbox-backend/pkg/common/email"
	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
	cMqtt "github.com/jovanfrandika/smartbox-backend/pkg/common/mqtt"
	cStorage "github.com/jovanfrandika/smartbox-backend/pkg/common/storage"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	deviceService "github.com/jovanfrandika/smartbox-backend/pkg/device/service"
	rFriendshipMongo "github.com/jovanfrandika/smartbox-backend/pkg/friendship/repository/mongo"
	friendshipService "github.com/jovanfrandika/smartbox-backend/pkg/friendship/service"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	rParcelRedis "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/redis"
	parcelService "github.com/jovanfrandika/smartbox-backend/pkg/parcel/service"
	rParcelTravelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/repository/mongo"
	parcelTravelService "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/service"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
	userService "github.com/jovanfrandika/smartbox-backend/pkg/user/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

func main() {
	log.Init()

	config.Init()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.Cfg.MongoDBUri))
	if err != nil {
		log.Fatal(err.Error(), 0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error(), 0)
	}

	db := client.Database(config.Cfg.DBName)

	userDb := rUserMongo.New(db)
	friendshipDb := rFriendshipMongo.New(db)
	parcelDb := rParcelMongo.New(db)
	parcelTravelDb := rParcelTravelMongo.New(db)
	deviceDb := rDeviceMongo.New(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Cfg.RedisAddr,
		Password: config.Cfg.RedisPassword,
		DB:       0,
	})

	parcelCache := rParcelRedis.New(rdb)

	gcsClient, err := gcs.NewClient(context.TODO(), option.WithCredentialsFile("files/service-account.json"))
	if err != nil {
		log.Fatal(err.Error(), 0)
	}
	defer gcsClient.Close()

	storage := cStorage.New(gcsClient, config.Cfg)

	email := cEmail.New(config.Cfg)

	mqttClient := cMqtt.Init("stancyzk", *config.Cfg)
	deviceMq := rDeviceMqtt.New(&mqttClient)

	r := chi.NewRouter()

	userRouter := chi.NewRouter()
	userService.Init(&userDb, userRouter, config.Cfg)

	deviceRouter := chi.NewRouter()
	deviceService.Init(&deviceDb, deviceRouter, config.Cfg)

	parcelRouter := chi.NewRouter()
	parcelService.Init(parcelService.InitInput{
		ParcelDb:    &parcelDb,
		ParcelCache: &parcelCache,
		UserDb:      &userDb,
		DeviceDb:    &deviceDb,
		DeviceMq:    &deviceMq,
		Config:      config.Cfg,
		Router:      parcelRouter,
		Storage:     &storage,
		Email:       &email,
	})

	parcelTravelRouter := chi.NewRouter()
	parcelTravelService.Init(&parcelTravelDb, &parcelDb, &deviceDb, parcelTravelRouter, config.Cfg)

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
	r.Mount("/parcel_travel", parcelTravelRouter)

	http.ListenAndServe(":8000", r)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
