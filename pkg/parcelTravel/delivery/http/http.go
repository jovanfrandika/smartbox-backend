package http

import (
	"context"
	"encoding/json"
	h "net/http"
	"time"

	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
	commonModel "github.com/jovanfrandika/smartbox-backend/pkg/common/model"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
)

const (
	TIMEOUT = 1 * time.Second
)

func (d *delivery) GetAll(w h.ResponseWriter, r *h.Request) {
	var payload model.GetAllInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetAllResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetAll(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Timeout", 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusRequestTimeout)
		json.NewEncoder(w).Encode(commonModel.ErrorResponse{
			Error: commonModel.TIMEOUT_ERROR,
		})
		return
	case <-ch:
		if err != nil {
			log.Error(err.Error(), 0)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(h.StatusInternalServerError)
			json.NewEncoder(w).Encode(commonModel.ErrorResponse{
				Error: commonModel.INTERVAL_SERVER_ERROR,
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
