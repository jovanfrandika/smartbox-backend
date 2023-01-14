package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) ConsumeTravelLog(_ m.Client, msg m.Message) {
	var payload model.ConsumeTravelLogMessage
	err := json.Unmarshal([]byte(msg.Payload()), &payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		msg.Ack()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.ConsumeTravelLog(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Create one parcel travel timeout")
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Create one parcel travel failed, Error: %v", err))
		}
		log.Info("Create one parcel travel success")
		msg.Ack()
	}
}
