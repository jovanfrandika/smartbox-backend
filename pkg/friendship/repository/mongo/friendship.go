package rMongo

import (
	"context"
	"errors"

	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Friendship struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   primitive.ObjectID `bson:"user_id,omitempty"`
	FriendID primitive.ObjectID `bson:"friend_user_id,omitempty"`
}

const (
	idField           = "_id"
	userIdField       = "user_id"
	friendUserIdField = "friend_user_id"
)

func (r *mongoDb) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) (string, error) {
	userID, err := primitive.ObjectIDFromHex(createOneInput.UserID)
	if err != nil {
		return "", err
	}

	friendUserId, err := primitive.ObjectIDFromHex(createOneInput.FriendUserID)
	if err != nil {
		return "", err
	}

	doc := bson.D{
		primitive.E{Key: userIdField, Value: userID},
		primitive.E{Key: friendUserIdField, Value: friendUserId},
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

func (r *mongoDb) DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error {
	userID, err := primitive.ObjectIDFromHex(deleteOneInput.UserID)
	if err != nil {
		return err
	}

	friendUserID, err := primitive.ObjectIDFromHex(deleteOneInput.FriendUserID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{
					"$and",
					bson.A{
						bson.D{{userIdField, bson.D{{"$eq", userID}}}},
						bson.D{{friendUserIdField, bson.D{{"$eq", friendUserID}}}},
					},
				}},
				bson.D{{
					"$and",
					bson.A{
						bson.D{{userIdField, bson.D{{"$eq", friendUserID}}}},
						bson.D{{friendUserIdField, bson.D{{"$eq", userID}}}},
					},
				}},
			},
		},
	}

	_, err = r.DbCollection.DeleteMany(ctx, filter, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoDb) GetAll(ctx context.Context, id string) ([]model.Friendship, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return []model.Friendship{}, err
	}

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{userIdField, bson.D{{"$eq", userID}}}},
				bson.D{{friendUserIdField, bson.D{{"$eq", userID}}}},
			},
		},
	}
	cursor, err := r.DbCollection.Find(ctx, filter, nil)
	if err != nil {
		return []model.Friendship{}, err
	}
	defer cursor.Close(nil)

	var output []model.Friendship
	for cursor.Next(ctx) {
		var elem Friendship
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.Friendship{}, err
		}
		output = append(output, model.Friendship{
			ID:           elem.ID.Hex(),
			UserID:       elem.UserID.Hex(),
			FriendUserID: elem.FriendID.Hex(),
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.Friendship{}, err
	}

	return output, nil
}
