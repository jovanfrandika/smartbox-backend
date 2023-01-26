package rMongo

import (
	"context"
	"errors"
	"time"

	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type ParcelTravel struct {
	ID         primitive.ObjectID `bson:"_id"`
	ParcelID   primitive.ObjectID `bson:"parcel_id"`
	Loc        Location           `bson:"loc"`
	Temp       float32            `bson:"temp"`
	Hmd        float32            `bson:"hmd"`
	DoorStatus int                `bson:"door_status"`
	Sgnl       int                `bson:"sgnl"`
	Spd        int                `bson:"spd"`
	Sats       int                `bson:"sats"`
	GPSTs      primitive.DateTime `bson:"gps_ts"`
	Ts         primitive.DateTime `bson:"ts"`
}

const (
	idField         = "_id"
	parcelIdField   = "parcel_id"
	locField        = "loc"
	tempField       = "temp"
	hmdField        = "hmd"
	doorStatusField = "door_status"
	sgnlField       = "sgnl"
	stlsField       = "stls"
	gpsTsField      = "gps_ts"
	tsField         = "ts"

	gpsTimeFormat = "2006-01-02T15:04:05Z"
)

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	parcelID, err := primitive.ObjectIDFromHex(createOneInput.ParcelID)
	if err != nil {
		return err
	}

	timestamp := time.Now()
	gpsTimestamp, err := time.Parse(gpsTimeFormat, createOneInput.GPSTs)
	if err != nil {
		return err
	}

	doc := bson.D{
		primitive.E{Key: parcelIdField, Value: parcelID},
		primitive.E{Key: locField, Value: bson.D{
			{Key: "type", Value: "Point"},
			{Key: "coordinates", Value: []float64{createOneInput.Coor.Lng, createOneInput.Coor.Lat}},
		}},
		primitive.E{Key: tempField, Value: createOneInput.Temp},
		primitive.E{Key: hmdField, Value: createOneInput.Hmd},
		primitive.E{Key: doorStatusField, Value: createOneInput.DoorStatus},
		primitive.E{Key: sgnlField, Value: createOneInput.Sgnl},
		primitive.E{Key: stlsField, Value: createOneInput.Stls},
		primitive.E{Key: gpsTsField, Value: primitive.NewDateTimeFromTime(gpsTimestamp.UTC())},
		primitive.E{Key: tsField, Value: primitive.NewDateTimeFromTime(timestamp.UTC())},
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
		{Key: parcelIdField, Value: parcelID},
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
			ID:       elem.ID.Hex(),
			ParcelID: elem.ParcelID.Hex(),
			Coor: model.Coordinate{
				Lat: elem.Loc.Coordinates[1],
				Lng: elem.Loc.Coordinates[0],
			},
			Temp:       elem.Temp,
			Hmd:        elem.Hmd,
			DoorStatus: elem.DoorStatus,
			Sgnl:       elem.Sgnl,
			GPSTs:      elem.GPSTs.Time(),
			Ts:         elem.Ts.Time(),
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.ParcelTravel{}, err
	}

	return output, nil
}
