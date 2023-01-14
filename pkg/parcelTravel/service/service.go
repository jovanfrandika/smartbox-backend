package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	dHttp "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/delivery/http"
	rParcelTravelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/usecase"
)

func Init(parcelTravelDb *rParcelTravelMongo.MongoDb, parcelDb *rParcelMongo.MongoDb, deviceDb *rDeviceMongo.MongoDb, r *chi.Mux, cfg *config.Config) {
	parcelTravelUsecase := usecase.New(cfg, parcelTravelDb, parcelDb, deviceDb)

	dHttp.Deliver(r, parcelTravelUsecase)
}
