package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/device/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(db *mongo.Database, r *chi.Mux, cfg *config.Config) {
	mongoRepo := rMongo.New(db)
	deviceUsecase := usecase.New(cfg, &mongoRepo)

	http.Deliver(r, deviceUsecase)
}
