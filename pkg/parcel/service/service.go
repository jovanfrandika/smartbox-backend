package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/email"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/storage"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/delivery/http"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	rParcelRedis "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/redis"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/usecase"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type InitInput struct {
	Config      *config.Config
	ParcelDb    *rParcelMongo.MongoDb
	ParcelCache *rParcelRedis.Redis
	UserDb      *rUserMongo.MongoDb
	DeviceDb    *rDeviceMongo.MongoDb
	DeviceMq    *rDeviceMqtt.Mqtt
	Storage     *storage.Storage
	Email       *email.Email
	Router      *chi.Mux
}

func Init(initInput InitInput) {
	parcelUsecase := usecase.New(usecase.NewInput{
		Config:      initInput.Config,
		ParcelDb:    initInput.ParcelDb,
		ParcelCache: initInput.ParcelCache,
		UserDb:      initInput.UserDb,
		DeviceDb:    initInput.DeviceDb,
		DeviceMq:    initInput.DeviceMq,
		Storage:     initInput.Storage,
		Email:       initInput.Email,
	})

	http.Deliver(initInput.Router, parcelUsecase)
}
