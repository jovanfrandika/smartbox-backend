package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/jwt"
	u "github.com/jovanfrandika/smartbox-backend/pkg/parcel/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(r *chi.Mux, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*r).With(jwt.AuthMiddleware).Get("/", d.Histories)
	(*r).With(jwt.AuthMiddleware).Post("/", d.CreateOne)
	(*r).With(jwt.AuthMiddleware).Put("/", d.UpdateOne)
	(*r).With(jwt.AuthMiddleware).Delete("/", d.DeleteOne)
	(*r).With(jwt.AuthMiddleware).Post("/photo/url", d.GetPhotoSignedUrl)
	(*r).With(jwt.AuthMiddleware).Post("/photo/check", d.CheckPhoto)
	(*r).With(jwt.AuthMiddleware).Post("/code/send", d.SendParcelCodeToReceiver)
	(*r).With(jwt.AuthMiddleware).Post("/code/verify", d.VerifyParcelCode)
	(*r).With(jwt.AuthMiddleware).Post("/progress", d.UpdateProgress)
	(*r).With(jwt.AuthMiddleware).Post("/open", d.OpenDoor)
}
