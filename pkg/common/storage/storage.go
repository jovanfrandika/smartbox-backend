package storage

import (
	"context"
	"time"

	gcs "cloud.google.com/go/storage"
	log "github.com/jovanfrandika/smartbox-backend/pkg/common/logger"
)

func (s *storage) GetSignedUrl(object string) (string, error) {
	opts := &gcs.SignedURLOptions{
		Scheme: gcs.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type:application/octet-stream",
		},
		Expires: time.Now().Add(15 * time.Minute),
	}
	url, err := s.client.Bucket(s.config.BucketName).SignedURL(object, opts)
	if err != nil {
		log.Error("Failed bucket object url", 0)
		return "", err
	}

	return url, nil
}

func (s *storage) IsObjectValid(ctx context.Context, object string) bool {
	obj := s.client.Bucket(s.config.BucketName).Object(object)
	_, err := obj.Attrs(ctx)
	if err != nil {
		log.Error("Invalid bucket object", 0)
		return false
	}

	return true
}
