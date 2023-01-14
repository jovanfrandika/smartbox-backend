package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/jwt"
	u "github.com/jovanfrandika/smartbox-backend/pkg/user/usecase"
)

type delivery struct {
	usecase u.Usecase
}

func Deliver(r *chi.Mux, usecase u.Usecase) {
	d := &delivery{
		usecase: usecase,
	}

	(*r).With(jwt.AuthMiddleware).Get("/me", d.Me)
	(*r).Post("/login", d.Login)
	(*r).Post("/register", d.Register)
	(*r).Post("/refresh", d.RefreshAccessToken)
}
