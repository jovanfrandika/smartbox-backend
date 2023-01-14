package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/jwt"
	u "github.com/jovanfrandika/smartbox-backend/pkg/device/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(r *chi.Mux, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*r).With(jwt.AuthMiddleware).Get("/", d.GetAll)
	(*r).With(jwt.AuthMiddleware).Post("/", d.CreateOne)
}
