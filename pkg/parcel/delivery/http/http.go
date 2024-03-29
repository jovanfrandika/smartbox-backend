package http

import (
	"context"
	"encoding/json"
	"fmt"
	h "net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
	commonModel "github.com/jovanfrandika/smartbox-backend/pkg/common/model"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
)

const (
	TIMEOUT = 2 * time.Second
)

func (d *delivery) GetOne(w h.ResponseWriter, r *h.Request) {
	var payload model.GetOneInput
	payload.ID = chi.URLParam(r, "ID")

	if payload.ID == "" {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetOneResponse
	var err error
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetOne(ctx, payload)
		ch <- 1
	}()

	select {
	case <-ctx.Done():
		log.Error("Timeout", 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusRequestTimeout)
		json.NewEncoder(w).Encode(commonModel.ErrorResponse{
			Error: commonModel.INTERVAL_SERVER_ERROR,
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
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) CreateOne(w h.ResponseWriter, r *h.Request) {
	var payload model.CreateOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		log.Error("Timeout", 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusRequestTimeout)
		json.NewEncoder(w).Encode(commonModel.ErrorResponse{
			Error: commonModel.INTERVAL_SERVER_ERROR,
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

func (d *delivery) UpdateOne(w h.ResponseWriter, r *h.Request) {
	var payload model.UpdateOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		log.Error("Timeout", 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusRequestTimeout)
		json.NewEncoder(w).Encode(commonModel.ErrorResponse{
			Error: commonModel.INTERVAL_SERVER_ERROR,
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
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) DeleteOne(w h.ResponseWriter, r *h.Request) {
	var payload model.DeleteOneInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		log.Error("Timeout", 0)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.StatusRequestTimeout)
		json.NewEncoder(w).Encode(commonModel.ErrorResponse{
			Error: commonModel.INTERVAL_SERVER_ERROR,
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
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) Histories(w h.ResponseWriter, r *h.Request) {
	var payload model.HistoryInput
	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
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

func (d *delivery) GetNearbyPickUps(w h.ResponseWriter, r *h.Request) {
	var payload model.GetNearbyPickUpsInput
	var err error
	payload.UserCoor.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}
	payload.UserCoor.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	var res model.GetNearbyPickUpsResponse
	ch := make(chan int)
	go func() {
		res, err = d.usecase.GetNearbyPickUps(ctx, payload)
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

func (d *delivery) GetPhotoSignedUrl(w h.ResponseWriter, r *h.Request) {
	var payload model.GetPhotoSignedUrlInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
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
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) CheckPhoto(w h.ResponseWriter, r *h.Request) {
	var payload model.CheckPhotoInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.CheckPhoto(ctx, payload)
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
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) SendParcelCode(w h.ResponseWriter, r *h.Request) {
	var payload model.SendParcelCodeInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}
	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.SendParcelCode(ctx, payload)
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
				Error: err.Error(),
			})
			return
		}
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) VerifyParcelCode(w h.ResponseWriter, r *h.Request) {
	var payload model.VerifyParcelCodeInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}
	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.VerifyParcelCode(ctx, payload)
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
				Error: err.Error(),
			})
			return
		}
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) UpdateProgress(w h.ResponseWriter, r *h.Request) {
	var payload model.UpdateProgressInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
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
		w.WriteHeader(h.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (d *delivery) OpenDoor(w h.ResponseWriter, r *h.Request) {
	var payload model.OpenDoorInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
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
		w.WriteHeader(h.StatusNoContent)
	}
}

func (d *delivery) CloseDoor(w h.ResponseWriter, r *h.Request) {
	var payload model.CloseDoorInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Error("Invalid Payload", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID")
	if reflect.TypeOf(userID).String() != "string" {
		log.Error("Invalid UserID", 0)
		w.WriteHeader(h.StatusBadRequest)
		return
	}

	payload.UserID = fmt.Sprintf("%v", userID)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ch := make(chan int)
	go func() {
		err = d.usecase.CloseDoor(ctx, payload)
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
		w.WriteHeader(h.StatusNoContent)
	}
}
