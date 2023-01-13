package rMqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mq struct {
	mqttClient *mqtt.Client
}

type Mqtt interface {
	PubOpenDoor(deviceName string) error
	PubStartTravel(deviceName string, lat float32, lng float32) error
	PubEndTravel(deviceName string) error
	PubStatus(deviceName string) error
}

func New(mqttClient *mqtt.Client) Mqtt {
	return &mq{
		mqttClient: mqttClient,
	}
}
