package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
)

func Init(deviceDb *rMongo.MongoDb, r *chi.Mux, cfg *config.Config) {
	deviceUsecase := usecase.New(cfg, deviceDb)

	http.Deliver(r, deviceUsecase)
}
