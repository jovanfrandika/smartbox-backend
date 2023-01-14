package rMongo

import (
	"context"
	"errors"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coordinate struct {
	Lat        float64 `bson:"lat"`
	Lng        float64 `bson:"lng"`
	Speed      int     `bson:"speed"`
	Satellites int     `bson:"satellites"`
}

type ParcelTravel struct {
	ID           primitive.ObjectID `bson:"_id"`
	ParcelID     primitive.ObjectID `bson:"parcel_id"`
	Coordinate   Coordinate         `bson:"coordinate"`
	IsDoorOpen   bool               `bson:"is_door_open"`
	Signal       int                `bson:"signal"`
	GPSTimestamp primitive.DateTime `bson:"gps_timestamp"`
	Timestamp    primitive.DateTime `bson:"timestamp"`
}

const (
	idField           = "_id"
	parcelIdField     = "parcel_id"
	coordinateField   = "coordinate"
	isDoorOpenField   = "is_door_open"
	signalField       = "signal"
	gpsTimestampField = "gps_timestamp"
	timestampField    = "timestamp"

	gpsTimeFormat = "2006-01-02T15:04:05Z"
)

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	parcelID, err := primitive.ObjectIDFromHex(createOneInput.ParcelID)
	if err != nil {
		return err
	}

	timestamp := time.Now()
	gpsTimestamp, err := time.Parse(gpsTimeFormat, createOneInput.GPSTimestamp)
	if err != nil {
		return err
	}

	doc := bson.D{
		primitive.E{Key: parcelIdField, Value: parcelID},
		primitive.E{Key: coordinateField, Value: Coordinate(createOneInput.Coordinate)},
		primitive.E{Key: isDoorOpenField, Value: createOneInput.IsDoorOpen},
		primitive.E{Key: signalField, Value: createOneInput.Signal},
		primitive.E{Key: gpsTimestampField, Value: primitive.NewDateTimeFromTime(gpsTimestamp.UTC())},
		primitive.E{Key: timestampField, Value: primitive.NewDateTimeFromTime(timestamp.UTC())},
	}

	res, err := r.DbCollection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); ok {
		return nil
	}

	return errors.New("No oid")
}

func (r *mongoDb) GetAll(ctx context.Context, getAllInput model.GetAllInput) ([]model.ParcelTravel, error) {
	parcelID, err := primitive.ObjectIDFromHex(getAllInput.ParcelID)
	if err != nil {
		return []model.ParcelTravel{}, err
	}

	filter := bson.D{
		{parcelIdField, parcelID},
	}
	cursor, err := r.DbCollection.Find(ctx, filter, nil)
	if err != nil {
		return []model.ParcelTravel{}, err
	}
	defer cursor.Close(nil)

	var output []model.ParcelTravel
	for cursor.Next(ctx) {
		var elem ParcelTravel
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.ParcelTravel{}, err
		}
		output = append(output, model.ParcelTravel{
			ID:           elem.ID.Hex(),
			ParcelID:     elem.ParcelID.Hex(),
			Coordinate:   model.Coordinate(elem.Coordinate),
			IsDoorOpen:   elem.IsDoorOpen,
			Signal:       elem.Signal,
			GPSTimestamp: elem.GPSTimestamp.Time(),
			Timestamp:    elem.Timestamp.Time(),
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.ParcelTravel{}, err
	}

	return output, nil
}
