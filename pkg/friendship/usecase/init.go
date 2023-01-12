package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
	rFriendshipDb "github.com/jovanfrandika/smartbox-backend/pkg/friendship/repository/mongo"
	rUserDb "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type usecase struct {
	config       *config.Config
	friendshipDb *rFriendshipDb.MongoDb
	userDb       *rUserDb.MongoDb
}

type Usecase interface {
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error
	GetAll(ctx context.Context, userID string) (model.GetAllResponse, error)
}

func New(config *config.Config, friendshipDb *rFriendshipDb.MongoDb, userDb *rUserDb.MongoDb) Usecase {
	return &usecase{
		config:       config,
		friendshipDb: friendshipDb,
		userDb:       userDb,
	}
}
