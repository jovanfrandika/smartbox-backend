package consumer

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	dMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/delivery/mqtt"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
)

func Init(deviceDb *rDeviceMongo.MongoDb, mqttClient *mqtt.Client, cfg *config.Config) {
	deviceUsecase := usecase.New(cfg, deviceDb)

	dMqtt.Deliver(mqttClient, deviceUsecase)
}
