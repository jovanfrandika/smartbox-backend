package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/utils"
	u "github.com/jovanfrandika/smartbox-backend/pkg/parcel/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(r *chi.Mux, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*r).With(utils.AuthMiddleware).Get("/", d.Histories)
	(*r).With(utils.AuthMiddleware).Post("/", d.CreateOne)
	(*r).With(utils.AuthMiddleware).Put("/", d.UpdateOne)
	(*r).With(utils.AuthMiddleware).Delete("/", d.DeleteOne)
	(*r).With(utils.AuthMiddleware).Get("/photo", d.GetPhotoSignedUrl)
	(*r).With(utils.AuthMiddleware).Post("/progress", d.UpdateProgress)
	(*r).With(utils.AuthMiddleware).Post("/open", d.OpenDoor)
}
