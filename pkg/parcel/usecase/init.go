package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/email"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/storage"
	rDeviceMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	rDeviceMqtt "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mqtt"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	rParcelMongo "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	rParcelRedis "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/redis"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type usecase struct {
	config      *config.Config
	parcelDb    *rParcelMongo.MongoDb
	parcelCache *rParcelRedis.Redis
	userDb      *rUserMongo.MongoDb
	deviceDb    *rDeviceMongo.MongoDb
	deviceMq    *rDeviceMqtt.Mqtt
	storage     *storage.Storage
	email       *email.Email
}

type NewInput struct {
	Config      *config.Config
	ParcelDb    *rParcelMongo.MongoDb
	ParcelCache *rParcelRedis.Redis
	UserDb      *rUserMongo.MongoDb
	DeviceDb    *rDeviceMongo.MongoDb
	DeviceMq    *rDeviceMqtt.Mqtt
	Storage     *storage.Storage
	Email       *email.Email
}

type Usecase interface {
	GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error)
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	Histories(ctx context.Context, historyInput model.HistoryInput) (model.HistoryResponse, error)
	GetNearbyPickUps(ctx context.Context, getNearyPickUpsInput model.GetNearbyPickUpsInput) (model.GetNearbyPickUpsResponse, error)
	GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error)
	CheckPhoto(ctx context.Context, checkPhotoInput model.CheckPhotoInput) error
	UpdateProgress(ctx context.Context, updateProgressInput model.UpdateProgressInput) (model.UpdateProgressResponse, error)
	OpenDoor(ctx context.Context, openDoorInput model.OpenDoorInput) error
	CloseDoor(ctx context.Context, closeDoorInput model.CloseDoorInput) error
	SendParcelCode(ctx context.Context, sendParcelCodeToReceiverInput model.SendParcelCodeInput) error
	VerifyParcelCode(ctx context.Context, verifyParcelCodeInput model.VerifyParcelCodeInput) error
}

func New(newInput NewInput) Usecase {
	return &usecase{
		config:      newInput.Config,
		parcelDb:    newInput.ParcelDb,
		parcelCache: newInput.ParcelCache,
		userDb:      newInput.UserDb,
		deviceDb:    newInput.DeviceDb,
		deviceMq:    newInput.DeviceMq,
		storage:     newInput.Storage,
		email:       newInput.Email,
	}
}
