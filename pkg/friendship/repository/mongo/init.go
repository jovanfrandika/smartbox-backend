package rMongo

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	DbCollection *mongo.Collection
}

type MongoDb interface {
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error)
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	GetAll(ctx context.Context, id string) ([]model.Friendship, error)
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
