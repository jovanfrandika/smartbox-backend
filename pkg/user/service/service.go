package userService

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/usecase"
)

func Init(userDb *rMongo.MongoDb, r *chi.Mux, cfg *config.Config) {
	userUsecase := usecase.New(cfg, userDb)

	http.Deliver(r, userUsecase)
}
