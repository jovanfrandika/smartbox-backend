package rMongo

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	DbCollection *mongo.Collection
}

type MongoDb interface {
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error)
	GetMany(ctx context.Context, deviceIDs []string) ([]model.Device, error)
	GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.Device, error)
	GetOneByName(ctx context.Context, getOneByNameInput model.GetOneByNameInput) (model.Device, error)
	GetAll(ctx context.Context) ([]model.Device, error)
	UpdateStatus(ctx context.Context, updateStatusInput model.UpdateStatusInput) error
}

const (
	collectionName = "device"
)

func New(dbClient *mongo.Database) MongoDb {
	collection := dbClient.Collection(collectionName)
	return &mongoDb{
		DbCollection: collection,
	}
}
