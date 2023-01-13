package usecase

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type usecase struct {
	config        *config.Config
	parcelDb      *rParcelMongo.MongoDb
	userDb        *rUserMongo.MongoDb
	deviceDb      *rDeviceMongo.MongoDb
	deviceMq      *rDeviceMqtt.Mqtt
	storageClient *storage.Client
}

type NewInput struct {
	Config        *config.Config
	ParcelDb      *rParcelMongo.MongoDb
	UserDb        *rUserMongo.MongoDb
	DeviceDb      *rDeviceMongo.MongoDb
	DeviceMq      *rDeviceMqtt.Mqtt
	StorageClient *storage.Client
}

type Usecase interface {
	GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error)
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	Histories(ctx context.Context, historyInput model.HistoryInput) (model.HistoryResponse, error)
	GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error)
	UpdateProgress(ctx context.Context, updateProgressInput model.UpdateProgressInput) (model.UpdateProgressResponse, error)
	OpenDoor(ctx context.Context, openDoorInput model.OpenDoorInput) error
}

func New(newInput NewInput) Usecase {
	return &usecase{
		config:        newInput.Config,
		parcelDb:      newInput.ParcelDb,
		userDb:        newInput.UserDb,
		deviceDb:      newInput.DeviceDb,
		deviceMq:      newInput.DeviceMq,
		storageClient: newInput.StorageClient,
	}
}
