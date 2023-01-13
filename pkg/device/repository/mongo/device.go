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

	IDLE_STATUS   = 0
	TRAVEL_STATUS = 1
)

func (r *mongoDb) GetMany(ctx context.Context, deviceIDs []string) ([]model.Device, error) {
	deviceObjectIDs := []primitive.ObjectID{}
	for _, deviceID := range deviceIDs {
		deviceObjectID, err := primitive.ObjectIDFromHex(deviceID)
		if err != nil {
			return []model.Device{}, err
		}
		deviceObjectIDs = append(deviceObjectIDs, deviceObjectID)
	}

	cursor, err := r.DbCollection.Find(ctx, bson.M{idField: bson.M{"$in": deviceObjectIDs}}, nil)
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

func (r *mongoDb) GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.Device, error) {
	docID, err := primitive.ObjectIDFromHex(getOneInput.ID)
	if err != nil {
		return model.Device{}, err
	}

	var res Device
	err = r.DbCollection.FindOne(ctx, bson.M{idField: docID}).Decode(&res)
	if err != nil {
		return model.Device{}, err
	}

	return model.Device{
		ID:     res.ID.Hex(),
		Name:   res.Name,
		Status: res.Status,
	}, nil
}

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error) {
	doc := bson.D{
		primitive.E{Key: nameField, Value: createOneInput.Name},
		primitive.E{Key: statusField, Value: IDLE_STATUS},
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

func (r *mongoDb) UpdateStatus(ctx context.Context, updateStatusInput model.UpdateStatusInput) error {
	docID, err := primitive.ObjectIDFromHex(updateStatusInput.ID)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: idField, Value: docID}}
	update := bson.D{
		primitive.E{Key: statusField, Value: IDLE_STATUS},
	}

	res := r.DbCollection.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (r *mongoDb) GetAll(ctx context.Context) ([]model.Device, error) {
	cursor, err := r.DbCollection.Find(ctx, bson.M{}, nil)
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
