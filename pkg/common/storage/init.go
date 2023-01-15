package storage

import (
	"context"

	gcs "cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
)

type storage struct {
	client *gcs.Client
	config *config.Config
}

type Storage interface {
	GetSignedUrl(object string) (string, error)
	IsObjectValid(ctx context.Context, object string) bool
}

func New(client *gcs.Client, config *config.Config) Storage {
	return &storage{
		client: client,
		config: config,
	}
}
