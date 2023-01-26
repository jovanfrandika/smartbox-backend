package rMqtt

import (
	"encoding/json"
)

const (
	CMD_TOPIC = "smartbox/cmd"

	OPEN_DOOR_DURATION  = "30"
	CLOSE_DOOR_DURATION = "0"

	CMD_OPEN_DOOR    = "door"
	CMD_START_TRAVEL = "start"
	CMD_END_TRAVEL   = "end"
	CMD_STATUS       = "status"
)

type Message struct {
	Cmd        string      `json:"cmd"`
	Value      interface{} `json:"value"`
	DeviceName string      `json:"device_name"`
}

func (r *mq) PubOpenDoor(deviceName string) error {
	message := Message{
		Cmd:        CMD_OPEN_DOOR,
		Value:      OPEN_DOOR_DURATION,
		DeviceName: deviceName,
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t := (*r.mqttClient).Publish(CMD_TOPIC, 1, false, []byte(messageStr))
	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (r *mq) PubCloseDoor(deviceName string) error {
	message := Message{
		Cmd:        CMD_OPEN_DOOR,
		Value:      CLOSE_DOOR_DURATION,
		DeviceName: deviceName,
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t := (*r.mqttClient).Publish(CMD_TOPIC, 1, false, []byte(messageStr))
	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (r *mq) PubStartTravel(deviceName string) error {
	message := Message{
		Cmd:        CMD_START_TRAVEL,
		Value:      "",
		DeviceName: deviceName,
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t := (*r.mqttClient).Publish(CMD_TOPIC, 1, false, []byte(messageStr))
	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (r *mq) PubEndTravel(deviceName string) error {
	message := Message{
		Cmd:        CMD_END_TRAVEL,
		Value:      nil,
		DeviceName: deviceName,
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t := (*r.mqttClient).Publish(CMD_TOPIC, 1, false, []byte(messageStr))
	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (r *mq) PubStatus(deviceName string) error {
	message := Message{
		Cmd:        CMD_STATUS,
		Value:      nil,
		DeviceName: deviceName,
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t := (*r.mqttClient).Publish(CMD_TOPIC, 1, false, []byte(messageStr))
	<-t.Done()

	if t.Error() != nil {
		return t.Error()
	}

	return nil
}
