package service

import (
	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/delivery/http"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/usecase"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type InitInput struct {
	Config        *config.Config
	ParcelDb      *rParcelMongo.MongoDb
	UserDb        *rUserMongo.MongoDb
	DeviceDb      *rDeviceMongo.MongoDb
	DeviceMq      *rDeviceMqtt.Mqtt
	StorageClient *storage.Client
	Router        *chi.Mux
}

func Init(initInput InitInput) {
	parcelUsecase := usecase.New(usecase.NewInput{
		Config:        initInput.Config,
		ParcelDb:      initInput.ParcelDb,
		UserDb:        initInput.UserDb,
		DeviceDb:      initInput.DeviceDb,
		DeviceMq:      initInput.DeviceMq,
		StorageClient: initInput.StorageClient,
	})

	http.Deliver(initInput.Router, parcelUsecase)
}
