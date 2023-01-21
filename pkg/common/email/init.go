package email

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
)

type email struct {
	config *config.Config
}

type Email interface {
	Send(ctx context.Context, sendInput SendInput) error
}

func New(config *config.Config) Email {
	return &email{
		config: config,
	}
}
