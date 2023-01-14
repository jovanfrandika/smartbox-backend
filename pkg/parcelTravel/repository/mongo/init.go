package rMongo

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	DbCollection *mongo.Collection
}

type MongoDb interface {
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	GetAll(ctx context.Context, getAllInput model.GetAllInput) ([]model.ParcelTravel, error)
}

const (
	collectionName = "parcel_travel"
)

func New(dbClient *mongo.Database) MongoDb {
	collection := dbClient.Collection(collectionName)
	return &mongoDb{
		DbCollection: collection,
	}
}
