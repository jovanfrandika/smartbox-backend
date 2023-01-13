package usecase

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type usecase struct {
	config        *config.Config
	parcelDb      *rParcelMongo.MongoDb
	userDb        *rUserMongo.MongoDb
	deviceDb      *rDeviceMongo.MongoDb
	storageClient *storage.Client
}

type Usecase interface {
	GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error)
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	Histories(ctx context.Context, historyInput model.HistoryInput) (model.HistoryResponse, error)
	GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error)
}

func New(config *config.Config, parcelDb *rParcelMongo.MongoDb, userDb *rUserMongo.MongoDb, deviceDb *rDeviceMongo.MongoDb, storageClient *storage.Client) Usecase {
	return &usecase{
		config:        config,
		parcelDb:      parcelDb,
		userDb:        userDb,
		deviceDb:      deviceDb,
		storageClient: storageClient,
	}
}
