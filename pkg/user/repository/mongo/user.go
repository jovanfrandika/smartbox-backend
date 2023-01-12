package rMongo

import (
	"context"
	"errors"

	"github.com/jovanfrandika/smartbox-backend/pkg/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
	Role     int                `bson:"role,omitempty"`
}

const (
	idField       = "_id"
	nameField     = "name"
	emailField    = "email"
	passwordField = "password"
	roleField     = "role"
	hashCost      = 14
)

func (r *mongoDb) CreateUser(ctx context.Context, registerInput model.RegisterInput) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), hashCost)
	if err != nil {
		return "", err
	}

	doc := bson.D{
		primitive.E{Key: nameField, Value: registerInput.Name},
		primitive.E{Key: emailField, Value: registerInput.Email},
		primitive.E{Key: passwordField, Value: string(hashedPassword)},
		primitive.E{Key: roleField, Value: registerInput.Role},
	}

	res, err := r.DbCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", errors.New("No oid")
}

func (r *mongoDb) GetUser(ctx context.Context, id string) (model.User, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	var res User
	err = r.DbCollection.FindOne(ctx, bson.M{idField: docID}).Decode(&res)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:    res.ID.Hex(),
		Name:  res.Name,
		Email: res.Email,
		Role:  res.Role,
	}, nil
}

func (r *mongoDb) GetMany(ctx context.Context, userIDs []string) ([]model.User, error) {
	userObjectIDs := []primitive.ObjectID{}
	for _, userID := range userIDs {
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return []model.User{}, err
		}
		userObjectIDs = append(userObjectIDs, userObjectID)
	}

	cursor, err := r.DbCollection.Find(ctx, bson.M{idField: bson.M{"$in": userObjectIDs}}, nil)
	if err != nil {
		return []model.User{}, err
	}
	defer cursor.Close(nil)

	var output []model.User
	for cursor.Next(ctx) {
		var elem User
		err := cursor.Decode(&elem)
		if err != nil {
			return []model.User{}, err
		}

		output = append(output, model.User{
			ID:    elem.ID.Hex(),
			Name:  elem.Name,
			Email: elem.Email,
			Role:  elem.Role,
		})
	}

	if err := cursor.Err(); err != nil {
		return []model.User{}, err
	}

	return output, nil
}

func (r *mongoDb) Login(ctx context.Context, loginInput model.LoginInput) (model.User, error) {
	var res User
	err := r.DbCollection.FindOne(ctx, bson.M{emailField: loginInput.Email}).Decode(&res)
	if err != nil {
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(loginInput.Password))
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:    res.ID.Hex(),
		Name:  res.Name,
		Email: res.Email,
		Role:  res.Role,
	}, nil
}
