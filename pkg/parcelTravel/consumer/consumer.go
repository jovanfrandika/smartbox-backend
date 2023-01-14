package consumer

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	mqtt "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/delivery/mqtt"
	rParcelTravelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/usecase"
)

func Init(parcelTravelDb *rParcelTravelMongo.MongoDb, parcelDb *rParcelMongo.MongoDb, deviceDb *rDeviceMongo.MongoDb, mqttClient *m.Client, cfg *config.Config) {
	parcelTravelUsecase := usecase.New(cfg, parcelTravelDb, parcelDb, deviceDb)

	mqtt.Deliver(mqttClient, parcelTravelUsecase)
}
