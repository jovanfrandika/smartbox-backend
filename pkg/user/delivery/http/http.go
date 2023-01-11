package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"reflect"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/user/model"
	log "github.com/sirupsen/logrus"
)

const (
	TIMEOUT = 5 * time.Second
)

func (d *delivery) Me(w h.ResponseWriter, r *h.Request) {
	var payload model.MeInput
	userId := r.Context().Value("userId")
	if reflect.TypeOf(userId).String() != "string" {
		log.Error("Error: Invalid UserId")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.ID = fmt.Sprintf("%v", userId)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.MeResponse
	var err error
	ch := make(chan int)
	go func() {
		res, err = d.usecase.Me(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Me timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Me failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) Login(w h.ResponseWriter, r *h.Request) {
	var payload model.LoginInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.LoginResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.Login(ctx, payload)
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

func (d *delivery) Register(w h.ResponseWriter, r *h.Request) {
	var payload model.RegisterInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.RegisterResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.Register(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Register timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Register failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) RefreshAccessToken(w h.ResponseWriter, r *h.Request) {
	var payload model.RefreshInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Error: Invalid Payload")
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.RefreshResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.RefreshAccessToken(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Refresh Access Token timeout")
		h.Error(w, "timeout", h.StatusInternalServerError)
		return
	case <-ch:
		if err != nil {
			log.Error(fmt.Sprintf("Refresh Access Token failed, Error: %v", err))
			h.Error(w, err.Error(), h.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
