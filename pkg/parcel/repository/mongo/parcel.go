package rMongo

import (
	"context"
	"errors"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type Photo struct {
	UpdatedAt primitive.DateTime `bson:"update_at"`
}

type Threshold struct {
	Low  float32 `bson:"low"`
	High float32 `bson:"high"`
}

type Parcel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Description  string             `bson:"description"`
	PickUpLoc    *Location          `bson:"pick_up_loc"`
	ArrivedLoc   *Location          `bson:"arrived_loc"`
	PickUpPhoto  *Photo             `bson:"pick_up_photo"`
	ArrivedPhoto *Photo             `bson:"arrived_photo"`
	TempThr      *Threshold         `bson:"temp_thr"`
	HmdThr       *Threshold         `bson:"hmd_thr"`
	ReceiverID   primitive.ObjectID `bson:"receiver_id"`
	SenderID     primitive.ObjectID `bson:"sender_id"`
	CourierID    primitive.ObjectID `bson:"courier_id"`
	DeviceID     primitive.ObjectID `bson:"device_id"`
	Status       int                `bson:"status"`
}

const (
	idField           = "_id"
	nameField         = "name"
	descriptionField  = "description"
	pickUpLocField    = "pick_up_loc"
	arrivedLocField   = "arrived_loc"
	pickUpPhotoField  = "pick_up_photo"
	arrivedPhotoField = "arrived_photo"
	uriField          = "uri"
	updatedAtField    = "updated_at"
	tempThrField      = "temp_thr"
	hmdThrField       = "hmd_thr"
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

	return buildParcel(res), nil
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

	return buildParcel(res), nil
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
		primitive.E{Key: pickUpLocField, Value: nil},
		primitive.E{Key: arrivedLocField, Value: nil},
		primitive.E{Key: pickUpPhotoField, Value: nil},
		primitive.E{Key: arrivedPhotoField, Value: nil},
		primitive.E{Key: tempThrField, Value: nil},
		primitive.E{Key: hmdThrField, Value: nil},
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

	var pickUpLoc *bson.D
	if updateOneInput.PickUpCoor != nil {
		pickUpLoc = &bson.D{
			{Key: "type", Value: "Point"},
			{Key: "coordinates", Value: []float64{updateOneInput.PickUpCoor.Lat, updateOneInput.PickUpCoor.Lng}},
		}
	}
	var arrivedLoc *bson.D
	if updateOneInput.ArrivedCoor != nil {
		arrivedLoc = &bson.D{
			{Key: "type", Value: "Point"},
			{Key: "coordinates", Value: []float64{updateOneInput.ArrivedCoor.Lat, updateOneInput.ArrivedCoor.Lng}},
		}
	}
	var pickUpPhoto *Photo
	if updateOneInput.PickUpPhoto != nil {
		pickUpPhoto = &Photo{
			UpdatedAt: primitive.NewDateTimeFromTime(updateOneInput.PickUpPhoto.UpdatedAt),
		}
	}
	var arrivedPhoto *Photo
	if updateOneInput.ArrivedPhoto != nil {
		arrivedPhoto = &Photo{
			UpdatedAt: primitive.NewDateTimeFromTime(updateOneInput.ArrivedPhoto.UpdatedAt),
		}
	}

	update := bson.D{
		primitive.E{Key: nameField, Value: updateOneInput.Name},
		primitive.E{Key: descriptionField, Value: updateOneInput.Description},
		primitive.E{Key: pickUpLocField, Value: pickUpLoc},
		primitive.E{Key: arrivedLocField, Value: arrivedLoc},
		primitive.E{Key: pickUpPhotoField, Value: pickUpPhoto},
		primitive.E{Key: arrivedPhotoField, Value: arrivedPhoto},
		primitive.E{Key: tempThrField, Value: updateOneInput.TempThr},
		primitive.E{Key: hmdThrField, Value: updateOneInput.HmdThr},
		primitive.E{Key: senderIdField, Value: senderID},
		primitive.E{Key: receiverIdField, Value: recieverID},
		primitive.E{Key: courierIdField, Value: courierID},
		primitive.E{Key: deviceIdField, Value: deviceID},
		primitive.E{Key: statusField, Value: updateOneInput.Status},
	}

	filter := bson.D{primitive.E{Key: idField, Value: docID}}

	_, err = r.DbCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
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
		{Key: "$or",
			Value: bson.A{
				bson.D{{Key: receiverIdField, Value: bson.D{{Key: "$eq", Value: userID}}}},
				bson.D{{Key: senderIdField, Value: bson.D{{Key: "$eq", Value: userID}}}},
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

		output = append(output, buildParcel(elem))
	}

	if err := cursor.Err(); err != nil {
		return []model.Parcel{}, err
	}

	return output, nil
}

func (r *mongoDb) GetNearbyPickUps(ctx context.Context, getNearbyPickUpsInput model.GetNearbyPickUpsInput) ([]model.Parcel, error) {
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: statusField, Value: bson.D{{Key: "$eq", Value: model.WAITING_FOR_COURIER_STATUS}}}},
				bson.D{{
					Key: pickUpLocField,
					Value: bson.D{{
						Key: "$near",
						Value: bson.D{
							{Key: "$geometry", Value: bson.D{
								{Key: "type", Value: "Point"},
								{Key: "coordinates", Value: []float64{getNearbyPickUpsInput.UserCoor.Lat, getNearbyPickUpsInput.UserCoor.Lng}},
							}},
							{Key: "$maxDistance", Value: 5000}, // in meters
						},
					}}}},
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

		output = append(output, buildParcel(elem))
	}

	if err := cursor.Err(); err != nil {
		return []model.Parcel{}, err
	}

	return output, nil
}
