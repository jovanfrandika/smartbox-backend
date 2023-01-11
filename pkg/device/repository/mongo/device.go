package rMongo

import (
	"context"
	"errors"

	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Device struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name,omitempty"`
	Status int                `bson:"status,omitempty"`
}

const (
	idField     = "_id"
	nameField   = "name"
	statusField = "status"
)

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error) {
	doc := bson.D{
		primitive.E{Key: nameField, Value: createOneInput.Name},
		primitive.E{Key: statusField, Value: createOneInput.Status},
	}

	res, err := r.DbCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.String(), nil
	}

	return "", errors.New("No oid")
}

func (r *mongoDb) GetAll(ctx context.Context) ([]model.Device, error) {
	cursor, err := r.DbCollection.Find(ctx, nil, nil)
	if err != nil {
		return []model.Device{}, err
	}
	defer cursor.Close(nil)

	var output []model.Device
	for cursor.Next(ctx) {
		var elem Device
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.Device{}, err
		}
		output = append(output, model.Device{
			ID:     elem.ID.Hex(),
			Name:   elem.Name,
			Status: elem.Status,
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.Device{}, err
	}

	return output, nil
}
