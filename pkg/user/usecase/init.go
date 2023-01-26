package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/model"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

type usecase struct {
	config *config.Config
	db     *rMongo.MongoDb
}

type Usecase interface {
	Me(ctx context.Context, registerInput model.MeInput) (model.MeResponse, error)
	Search(ctx context.Context, searchInput model.SearchInput) (model.SearchResponse, error)
	Register(ctx context.Context, registerInput model.RegisterInput) (model.RegisterResponse, error)
	Login(ctx context.Context, loginInput model.LoginInput) (model.LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshInput model.RefreshInput) (model.RefreshResponse, error)
}

func New(config *config.Config, db *rMongo.MongoDb) Usecase {
	return &usecase{
		config: config,
		db:     db,
	}
}
