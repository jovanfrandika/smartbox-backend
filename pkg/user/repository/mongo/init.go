package rMongo

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/user/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	DbCollection *mongo.Collection
}

type MongoDb interface {
	CreateUser(ctx context.Context, registerInput model.RegisterInput) (string, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	GetMany(ctx context.Context, userIDs []string) ([]model.User, error)
	Login(ctx context.Context, loginInput model.LoginInput) (model.User, error)
}

const (
	collectionName = "user"
)

func New(dbClient *mongo.Database) MongoDb {
	collection := dbClient.Collection(collectionName)
	return &mongoDb{
		DbCollection: collection,
	}
}
