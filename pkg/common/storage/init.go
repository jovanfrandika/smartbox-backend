package storage

import (
	"context"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
)

type storage struct {
	client *gcs.Client
	config *config.Config
}

type Storage interface {
	GetSignedUrl(object string) (string, error)
	GetObjectUpdatedAt(ctx context.Context, object string) (time.Time, error)
}

func New(client *gcs.Client, config *config.Config) Storage {
	return &storage{
		client: client,
		config: config,
	}
}
