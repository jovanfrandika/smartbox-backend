package rMongo

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/travel/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	DbCollection *mongo.Collection
}

type MongoDb interface {
	GetOne(ctx context.Context, id string) (model.Travel, error)
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error)
	UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	Histories(ctx context.Context, historyInput model.HistoryInput) ([]model.Travel, error)
}

const (
	collectionName = "travel"
)

func New(dbClient *mongo.Database) MongoDb {
	collection := dbClient.Collection(collectionName)
	return &mongoDb{
		DbCollection: collection,
	}
}
