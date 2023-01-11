package container

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/user/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(db *mongo.Database, r *chi.Mux, cfg *config.Config) {
	mongoRepo := rMongo.New(db)
	userUsecase := usecase.New(cfg, &mongoRepo)

	http.Deliver(r, userUsecase)
}
