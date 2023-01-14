package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceDb "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rParcelDb "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
	rParcelTravelDb "github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/repository/mongo"
)

type usecase struct {
	config         *config.Config
	parcelTravelDb *rParcelTravelDb.MongoDb
	parcelDb       *rParcelDb.MongoDb
	deviceDb       *rDeviceDb.MongoDb
}

type Usecase interface {
	ConsumeTravelLog(ctx context.Context, consumeTravelLogMessage model.ConsumeTravelLogMessage) error
	GetAll(ctx context.Context, getAllInput model.GetAllInput) (model.GetAllResponse, error)
}

func New(config *config.Config, parcelTravelDb *rParcelTravelDb.MongoDb, parcelDb *rParcelDb.MongoDb, deviceDb *rDeviceDb.MongoDb) Usecase {
	return &usecase{
		config:         config,
		parcelTravelDb: parcelTravelDb,
		parcelDb:       parcelDb,
		deviceDb:       deviceDb,
	}
}
