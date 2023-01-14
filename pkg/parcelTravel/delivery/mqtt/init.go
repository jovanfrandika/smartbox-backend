package mqtt

import (
	m "github.com/eclipse/paho.mqtt.golang"
	u "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/usecase"
)

type delivery struct {
	usecase u.Usecase
}

const (
	SMARTBOX_DATA_TOPIC = "smartbox/data"
)

func Deliver(mqttClient *m.Client, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*mqttClient).Subscribe(SMARTBOX_DATA_TOPIC, 1, d.ConsumeTravelLog)
}
