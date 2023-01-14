package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	u "github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
)

type delivery struct {
	usecase u.Usecase
}

const (
	SMARTBOX_STATUS_TOPIC = "smartbox/status"
)

func Deliver(mqttClient *mqtt.Client, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*mqttClient).Subscribe(SMARTBOX_STATUS_TOPIC, 1, d.ConsumeUpdateStatusLog)
}
