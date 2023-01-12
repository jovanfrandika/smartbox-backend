package service

import (
	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/travel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(db *mongo.Database, r *chi.Mux, cfg *config.Config, storageClient *storage.Client) {
	mongoRepo := rMongo.New(db)
	travelUsecase := usecase.New(cfg, &mongoRepo, storageClient)

	http.Deliver(r, travelUsecase)
}
