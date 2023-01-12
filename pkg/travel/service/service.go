package service

import (
	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/delivery/http"
	rMongo "github.com/jovanfrandika/smartbox-backend/pkg/travel/repository/mongo"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/usecase"
)

func Init(travelDb *rMongo.MongoDb, r *chi.Mux, cfg *config.Config, storageClient *storage.Client) {
	travelUsecase := usecase.New(cfg, travelDb, storageClient)

	http.Deliver(r, travelUsecase)
}
