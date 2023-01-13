package service

import (
	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/delivery/http"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/usecase"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

func Init(parcelDb *rParcelMongo.MongoDb, userDb *rUserMongo.MongoDb, deviceDb *rDeviceMongo.MongoDb, r *chi.Mux, cfg *config.Config, storageClient *storage.Client) {
	parcelUsecase := usecase.New(cfg, parcelDb, userDb, deviceDb, storageClient)

	http.Deliver(r, parcelUsecase)
}
