package usecase

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/model"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/travel/repository/mongo"
)

type usecase struct {
	config        *config.Config
	db            *rMongo.MongoDb
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

func New(config *config.Config, db *rMongo.MongoDb, storageClient *storage.Client) Usecase {
	return &usecase{
		config:        config,
		db:            db,
		storageClient: storageClient,
	}
}
