package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/utils"
	u "github.com/jovanfrandika/smartbox-backend/pkg/friendship/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(r *chi.Mux, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*r).With(utils.AuthMiddleware).Get("/", d.GetAll)
	(*r).With(utils.AuthMiddleware).Post("/", d.CreateOne)
	(*r).With(utils.AuthMiddleware).Delete("/", d.DeleteOne)
}
