package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) CreateOne(w h.ResponseWriter, r *h.Request) {
	var payload model.CreateOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.CreateOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Create one device timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Me failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusCreated)
	}
}

func (d *delivery) GetAll(w h.ResponseWriter, r *h.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetAllResponse
	var err error
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetAll(ctx)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Login timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Login failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
