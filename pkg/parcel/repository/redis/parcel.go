package rRedis

import (
	"context"
	"fmt"
	"time"
)

const (
	PARCEL_CODE_KEY = "parcel_code:%s"

	PARCEL_CODE_EXPIRATION = 1 * 60 * time.Second

	CODE_CHARS = "1234567890"
)

func (r *redis) SetParcelCode(ctx context.Context, parcel_id, code string) error {
	err := (*r.redisClient).Set(ctx, fmt.Sprintf(PARCEL_CODE_KEY, parcel_id), code, PARCEL_CODE_EXPIRATION).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) GetParcelCode(ctx context.Context, parcel_id string) (string, error) {
	res, err := (*r.redisClient).Get(ctx, fmt.Sprintf(PARCEL_CODE_KEY, parcel_id)).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
