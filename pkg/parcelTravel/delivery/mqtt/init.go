package mqtt

import (
	m "github.com/eclipse/paho.mqtt.golang"
	u "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(mqttClient *m.Client, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*mqttClient).AddRoute("smartbox/data", d.ConsumeTravelLog)
}
