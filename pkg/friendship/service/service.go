package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/delivery/http"
	rFriendshipMongo "github.com/jovanfrandika/smartbox-backend/pkg/friendship/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/usecase"
	rUserMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
)

func Init(friendshipDb *rFriendshipMongo.MongoDb, userDb *rUserMongo.MongoDb, r *chi.Mux, cfg *config.Config) {
	friendshipUsecase := usecase.New(cfg, friendshipDb, userDb)

	http.Deliver(r, friendshipUsecase)
}
