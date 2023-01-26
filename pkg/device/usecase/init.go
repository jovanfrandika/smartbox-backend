package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
)

type usecase struct {
	config *config.Config
	db     *rMongo.MongoDb
}

type Usecase interface {
	ConsumeUpdateStatus(ctx context.Context, consumeUpdateStatusInput model.ConsumeUpdateStatusMessage) error
	CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error
	GetAll(ctx context.Context) (model.GetAllResponse, error)
	GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error)
	GetQRCode(ctx context.Context, getQRCodeInput model.GetQRCodeInput) (model.GetQRCodeResponse, error)
}

func New(config *config.Config, db *rMongo.MongoDb) Usecase {
	return &usecase{
		config: config,
		db:     db,
	}
}
