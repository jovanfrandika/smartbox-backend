package rMongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jovanfrandika/smartbox-backend/pkg/travel/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coordinate struct {
	Lat  float32 `bson:"lat,omitempty"`
	Long float32 `bson:"long,omitempty"`
}

type Travel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	PhotoUri    string             `bson:"photo_uri,omitempty"`
	Start       Coordinate         `bson:"start,omitempty"`
	End         Coordinate         `bson:"end,omitempty"`
	ReceiverID  primitive.ObjectID `bson:"receiver_id,omitempty"`
	SenderID    primitive.ObjectID `bson:"sender_id,omitempty"`
	Status      int                `bson:"status,omitempty"`
}

const (
	idField          = "_id"
	nameField        = "name"
	descriptionField = "description"
	photoUriField    = "photo_uri"
	startField       = "start"
	endField         = "end"
	recieverIdField  = "reciever_id"
	senderIdField    = "sender_id"
	statusField      = "status"

	DRAFT_STATUS    = 0
	ON_GOING_STATUS = 1
	DONE_STATUS     = 2
)

func (r *mongoDb) GetOne(ctx context.Context, id string) (model.Travel, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Travel{}, err
	}

	var res Travel
	err = r.DbCollection.FindOne(ctx, bson.M{idField: docId}).Decode(&res)
	if err != nil {
		return model.Travel{}, err
	}

	return model.Travel{
		ID:          res.ID.Hex(),
		Name:        res.Name,
		Description: res.Description,
		PhotoUri:    res.PhotoUri,
		Start:       model.Coordinate(res.Start),
		End:         model.Coordinate(res.End),
		ReceiverID:  res.ReceiverID.Hex(),
		SenderID:    res.SenderID.Hex(),
		Status:      res.Status,
	}, nil
}

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error) {
	docID := primitive.NewObjectID()

	senderID, err := primitive.ObjectIDFromHex(createOneInput.SenderID)
	if err != nil {
		return "", err
	}

	doc := bson.D{
		primitive.E{Key: idField, Value: docID},
		primitive.E{Key: nameField, Value: ""},
		primitive.E{Key: descriptionField, Value: ""},
		primitive.E{Key: photoUriField, Value: fmt.Sprintf("%s/%s.jpg", createOneInput.SenderID, docID.Hex())},
		primitive.E{Key: startField, Value: nil},
		primitive.E{Key: endField, Value: nil},
		primitive.E{Key: recieverIdField, Value: nil},
		primitive.E{Key: senderIdField, Value: senderID},
		primitive.E{Key: statusField, Value: DRAFT_STATUS},
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

func (r *mongoDb) UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error {
	docID, err := primitive.ObjectIDFromHex(updateOneInput.ID)
	if err != nil {
		return err
	}

	senderID, err := primitive.ObjectIDFromHex(updateOneInput.SenderID)
	if err != nil {
		return err
	}

	update := bson.D{
		primitive.E{Key: nameField, Value: updateOneInput.Name},
		primitive.E{Key: descriptionField, Value: updateOneInput.Description},
		primitive.E{Key: photoUriField, Value: updateOneInput.PhotoUri},
		primitive.E{Key: startField, Value: updateOneInput.Start},
		primitive.E{Key: endField, Value: updateOneInput.End},
		primitive.E{Key: senderIdField, Value: senderID},
		primitive.E{Key: statusField, Value: DRAFT_STATUS},
	}

	if updateOneInput.ReceiverID != "" {
		recieverID, err := primitive.ObjectIDFromHex(updateOneInput.ReceiverID)
		if err != nil {
			return err
		}
		update = append(update, primitive.E{Key: recieverIdField, Value: recieverID})
	}

	filter := bson.D{primitive.E{Key: idField, Value: docID}}

	_, err = r.DbCollection.UpdateOne(ctx, filter, bson.D{{"$set", update}})
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoDb) DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error {
	docID, err := primitive.ObjectIDFromHex(deleteOneInput.ID)
	if err != nil {
		return err
	}

	_, err = r.DbCollection.DeleteOne(ctx, bson.M{idField: docID})
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoDb) Histories(ctx context.Context, historyInput model.HistoryInput) ([]model.Travel, error) {
	userID, err := primitive.ObjectIDFromHex(historyInput.UserID)
	if err != nil {
		return []model.Travel{}, err
	}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{recieverIdField, bson.D{{"$eq", userID}}}},
				bson.D{{senderIdField, bson.D{{"$eq", userID}}}},
			},
		},
	}
	cursor, err := r.DbCollection.Find(ctx, filter, nil)
	if err != nil {
		return []model.Travel{}, err
	}
	defer cursor.Close(nil)

	var output []model.Travel
	for cursor.Next(ctx) {
		var elem Travel
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.Travel{}, err
		}

		receiverID := ""
		if elem.ReceiverID != primitive.NilObjectID {
			receiverID = elem.ReceiverID.Hex()
		}

		senderID := ""
		if elem.SenderID != primitive.NilObjectID {
			senderID = elem.SenderID.Hex()
		}

		output = append(output, model.Travel{
			ID:          elem.ID.Hex(),
			Name:        elem.Name,
			Description: elem.Description,
			PhotoUri:    elem.PhotoUri,
			Start:       model.Coordinate(elem.Start),
			End:         model.Coordinate(elem.End),
			ReceiverID:  receiverID,
			SenderID:    senderID,
			Status:      elem.Status,
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.Travel{}, err
	}

	return output, nil
}
