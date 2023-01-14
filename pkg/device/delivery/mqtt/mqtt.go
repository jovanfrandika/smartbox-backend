package mqtt

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) ConsumeUpdateStatusLog(_ m.Client, msg m.Message) {
	var payload model.ConsumeUpdateStatusMessage
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		msg.Ack()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.ConsumeUpdateStatus(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Timeout", 0)
		return
	case <-ch:
		if err != nil {
			log.Error(err.Error(), 0)
			msg.Ack()
			return
		}
		log.Info("Update device status success", 0)
		msg.Ack()
	}
}
