package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"reflect"
	"time"

	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
	commonModel "github.com/jovanfrandika/smartbox-backend/pkg/common/model"

	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) CreateOne(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	var payload model.CreateOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		log.Error(err.Error(), 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusInternalServerError)
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
		w.WriteHeader(h.StatusCreated)
	}
}

func (d *delivery) DeleteOne(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	var payload model.DeleteOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		w.WriteHeader(h.StatusCreated)
	}
}

func (d *delivery) GetAll(w h.ResponseWriter, r *h.Request) {
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
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
