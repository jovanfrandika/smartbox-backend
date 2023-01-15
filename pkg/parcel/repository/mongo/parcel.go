package rMongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coordinate struct {
	Lat  float32 `bson:"lat,omitempty"`
	Long float32 `bson:"long,omitempty"`
}

type Parcel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name,omitempty"`
	Description  string             `bson:"description,omitempty"`
	PhotoUri     string             `bson:"photo_uri,omitempty"`
	IsPhotoValid bool               `bson:"is_photo_valid,omitempty"`
	Start        *Coordinate        `bson:"start,omitempty"`
	End          *Coordinate        `bson:"end,omitempty"`
	ReceiverID   primitive.ObjectID `bson:"receiver_id,omitempty"`
	SenderID     primitive.ObjectID `bson:"sender_id,omitempty"`
	CourierID    primitive.ObjectID `bson:"courier_id,omitempty"`
	DeviceID     primitive.ObjectID `bson:"device_id,omitempty"`
	Status       int                `bson:"status,omitempty"`
}

const (
	idField           = "_id"
	nameField         = "name"
	descriptionField  = "description"
	photoUriField     = "photo_uri"
	isPhotoValidField = "is_photo_valid"
	startField        = "start"
	endField          = "end"
	receiverIdField   = "receiver_id"
	senderIdField     = "sender_id"
	courierIdField    = "courier_id"
	deviceIdField     = "device_id"
	statusField       = "status"

	EmptyObjectId = "000000000000000000000000"
)

func (r *mongoDb) GetOne(ctx context.Context, id string) (model.Parcel, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Parcel{}, err
	}

	var res Parcel
	err = r.DbCollection.FindOne(ctx, bson.M{idField: docID}).Decode(&res)
	if err != nil {
		return model.Parcel{}, err
	}

	var start *model.Coordinate
	if res.Start != nil {
		start = &model.Coordinate{
			Lat:  res.Start.Lat,
			Long: res.Start.Long,
		}
	}
	var end *model.Coordinate
	if res.End != nil {
		end = &model.Coordinate{
			Lat:  res.End.Lat,
			Long: res.End.Long,
		}
	}

	return model.Parcel{
		ID:           res.ID.Hex(),
		Name:         res.Name,
		Description:  res.Description,
		PhotoUri:     res.PhotoUri,
		IsPhotoValid: res.IsPhotoValid,
		Start:        start,
		End:          end,
		ReceiverID:   res.ReceiverID.Hex(),
		SenderID:     res.SenderID.Hex(),
		CourierID:    res.CourierID.Hex(),
		DeviceID:     res.DeviceID.Hex(),
		Status:       res.Status,
	}, nil
}

func (r *mongoDb) GetOneByDeviceAndStatus(ctx context.Context, getOneByDeviceAndStatusInput model.GetOneByDeviceAndStatusInput) (model.Parcel, error) {
	deviceID, err := primitive.ObjectIDFromHex(getOneByDeviceAndStatusInput.Device)
	if err != nil {
		return model.Parcel{}, err
	}

	filter := bson.M{
		"$and": bson.A{
			bson.M{deviceIdField: deviceID},
			bson.M{statusField: getOneByDeviceAndStatusInput.Status},
		},
	}
	var res Parcel
	err = r.DbCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return model.Parcel{}, err
	}

	var start *model.Coordinate
	if res.Start != nil {
		start = &model.Coordinate{
			Lat:  res.Start.Lat,
			Long: res.Start.Long,
		}
	}
	var end *model.Coordinate
	if res.End != nil {
		end = &model.Coordinate{
			Lat:  res.End.Lat,
			Long: res.End.Long,
		}
	}

	return model.Parcel{
		ID:           res.ID.Hex(),
		Name:         res.Name,
		Description:  res.Description,
		PhotoUri:     res.PhotoUri,
		IsPhotoValid: res.IsPhotoValid,
		Start:        start,
		End:          end,
		ReceiverID:   res.ReceiverID.Hex(),
		SenderID:     res.SenderID.Hex(),
		CourierID:    res.CourierID.Hex(),
		DeviceID:     res.DeviceID.Hex(),
		Status:       res.Status,
	}, nil
}

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error) {
	docID := primitive.NewObjectID()

	senderID, err := primitive.ObjectIDFromHex(createOneInput.SenderID)
	if err != nil {
		return "", err
	}

	emptyID, err := primitive.ObjectIDFromHex(EmptyObjectId)
	if err != nil {
		return "", err
	}

	doc := bson.D{
		primitive.E{Key: idField, Value: docID},
		primitive.E{Key: nameField, Value: ""},
		primitive.E{Key: descriptionField, Value: ""},
		primitive.E{Key: photoUriField, Value: fmt.Sprintf("%s/%s.jpg", createOneInput.SenderID, docID.Hex())},
		primitive.E{Key: isPhotoValidField, Value: false},
		primitive.E{Key: startField, Value: nil},
		primitive.E{Key: endField, Value: nil},
		primitive.E{Key: receiverIdField, Value: emptyID},
		primitive.E{Key: senderIdField, Value: senderID},
		primitive.E{Key: courierIdField, Value: emptyID},
		primitive.E{Key: deviceIdField, Value: emptyID},
		primitive.E{Key: statusField, Value: model.DRAFT_STATUS},
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

	recieverID, err := primitive.ObjectIDFromHex(updateOneInput.ReceiverID)
	if err != nil {
		return err
	}

	courierID, err := primitive.ObjectIDFromHex(updateOneInput.CourierID)
	if err != nil {
		return err
	}

	deviceID, err := primitive.ObjectIDFromHex(updateOneInput.DeviceID)
	if err != nil {
		return err
	}

	update := bson.D{
		primitive.E{Key: nameField, Value: updateOneInput.Name},
		primitive.E{Key: descriptionField, Value: updateOneInput.Description},
		primitive.E{Key: photoUriField, Value: updateOneInput.PhotoUri},
		primitive.E{Key: isPhotoValidField, Value: updateOneInput.IsPhotoValid},
		primitive.E{Key: startField, Value: updateOneInput.Start},
		primitive.E{Key: endField, Value: updateOneInput.End},
		primitive.E{Key: senderIdField, Value: senderID},
		primitive.E{Key: receiverIdField, Value: recieverID},
		primitive.E{Key: courierIdField, Value: courierID},
		primitive.E{Key: deviceIdField, Value: deviceID},
		primitive.E{Key: statusField, Value: updateOneInput.Status},
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

func (r *mongoDb) Histories(ctx context.Context, historyInput model.HistoryInput) ([]model.Parcel, error) {
	userID, err := primitive.ObjectIDFromHex(historyInput.UserID)
	if err != nil {
		return []model.Parcel{}, err
	}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{receiverIdField, bson.D{{"$eq", userID}}}},
				bson.D{{senderIdField, bson.D{{"$eq", userID}}}},
			},
		},
	}
	cursor, err := r.DbCollection.Find(ctx, filter, nil)
	if err != nil {
		return []model.Parcel{}, err
	}
	defer cursor.Close(nil)

	var output []model.Parcel
	for cursor.Next(ctx) {
		var elem Parcel
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.Parcel{}, err
		}

		var start *model.Coordinate
		if elem.Start != nil {
			start = &model.Coordinate{
				Lat:  elem.Start.Lat,
				Long: elem.Start.Long,
			}
		}
		var end *model.Coordinate
		if elem.End != nil {
			end = &model.Coordinate{
				Lat:  elem.End.Lat,
				Long: elem.End.Long,
			}
		}

		output = append(output, model.Parcel{
			ID:           elem.ID.Hex(),
			Name:         elem.Name,
			Description:  elem.Description,
			PhotoUri:     elem.PhotoUri,
			IsPhotoValid: elem.IsPhotoValid,
			Start:        start,
			End:          end,
			ReceiverID:   elem.ReceiverID.Hex(),
			SenderID:     elem.SenderID.Hex(),
			CourierID:    elem.CourierID.Hex(),
			DeviceID:     elem.DeviceID.Hex(),
			Status:       elem.Status,
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.Parcel{}, err
	}

	return output, nil
}
