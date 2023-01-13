package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"reflect"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 2 * time.Second
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

func (d *delivery) UpdateOne(w h.ResponseWriter, r *h.Request) {
	var payload model.UpdateOneInput
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
		err = d.usecase.UpdateOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Update one device timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Me failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) DeleteOne(w h.ResponseWriter, r *h.Request) {
	var payload model.DeleteOneInput
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
		err = d.usecase.DeleteOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Delete one device timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Me failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) Histories(w h.ResponseWriter, r *h.Request) {
	var payload model.HistoryInput
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.HistoryResponse
	var err error
	ch := make(chan int)
	go func() {
		res, err = d.usecase.Histories(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Histories timeout")
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

func (d *delivery) GetPhotoSignedUrl(w h.ResponseWriter, r *h.Request) {
	var payload model.GetPhotoSignedUrlInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetPhotoSignedUrlResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetPhotoSignedUrl(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Delete one device timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Me failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) UpdateProgress(w h.ResponseWriter, r *h.Request) {
	var payload model.UpdateProgressInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.UpdateProgressResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.UpdateProgress(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Update progress device timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Update progress failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) OpenDoor(w h.ResponseWriter, r *h.Request) {
	var payload model.OpenDoorInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Error: Invalid UserID")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.OpenDoor(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Open Door timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Open Door failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.WriteHeader(h.StatusNoContent)
	}
}
