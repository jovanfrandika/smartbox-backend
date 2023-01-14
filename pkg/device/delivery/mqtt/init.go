package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	u "github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(mqttClient *mqtt.Client, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*mqttClient).AddRoute("smartbox/status", d.ConsumeUpdateStatusLog)
}
