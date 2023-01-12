package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"reflect"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) CreateOne(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	var payload model.CreateOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.CreateOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Create one friendship timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Create one friendship failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusCreated)
	}
}

func (d *delivery) DeleteOne(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	var payload model.DeleteOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.DeleteOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Delete one friendship timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Delete one failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusCreated)
	}
}

func (d *delivery) GetAll(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetAllResponse
	var err error
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetAll(ctx, fmt.Sprintf("%v", userID))
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Get all timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Get all failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
