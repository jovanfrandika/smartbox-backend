package rRedis

import (
	"context"

	r "github.com/go-redis/redis/v9"
)

type redis struct {
	redisClient *r.Client
}

type Redis interface {
	SetParcelCode(ctx context.Context, parcel_id, code string) error
	GetParcelCode(ctx context.Context, parcel_id string) (string, error)
}

func New(redisClient *r.Client) Redis {
	return &redis{
		redisClient: redisClient,
	}
}
