package rMqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mq struct {
	mqttClient *mqtt.Client
}

type Mqtt interface {
	PubOpenDoor(deviceName string) error
	PubCloseDoor(deviceName string) error
	PubStartTravel(deviceName string) error
	PubEndTravel(deviceName string) error
	PubStatus(deviceName string) error
}

func New(mqttClient *mqtt.Client) Mqtt {
	return &mq{
		mqttClient: mqttClient,
	}
}
