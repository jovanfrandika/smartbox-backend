package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jovanfrandika/smartbox-backend/pkg/travel/model"
	travelCol "github.com/jovanfrandika/smartbox-backend/pkg/travel/repository/mongo"
	log "github.com/sirupsen/logrus"
)

func (u *usecase) GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error) {
	res, err := (*u.db).GetOne(ctx, getOneInput.ID)
	if err != nil {
		return model.GetOneResponse{}, err
	}

	return model.GetOneResponse(res), nil
}

func (u *usecase) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	_, err := (*u.db).CreateOne(ctx, createOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error {
	err := (*u.db).UpdateOne(ctx, updateOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error {
	res, err := (*u.db).GetOne(ctx, deleteOneInput.ID)
	if err != nil {
		return err
	}

	if res.Status != travelCol.DRAFT_STATUS {
		return errors.New("Cannot delete non draft travel")
	}

	err = (*u.db).DeleteOne(ctx, deleteOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) Histories(ctx context.Context, historyInput model.HistoryInput) (model.HistoryResponse, error) {
	histories, err := (*u.db).Histories(ctx, historyInput)
	if err != nil {
		log.Error(fmt.Sprintf("Error: %s", err.Error()))
		return model.HistoryResponse{}, err
	}

	return model.HistoryResponse{
		Histories: histories,
	}, nil
}

func (u *usecase) GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error) {
	res, err := (*u.db).GetOne(ctx, getPhotoSignedUrlInput.ID)
	if err != nil {
		return model.GetPhotoSignedUrlResponse{}, err
	}

	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type:application/octet-stream",
		},
		Expires: time.Now().Add(15 * time.Minute),
	}
	url, err := u.storageClient.Bucket(u.config.BucketName).SignedURL(res.PhotoUri, opts)
	if err != nil {
		return model.GetPhotoSignedUrlResponse{}, err
	}

	return model.GetPhotoSignedUrlResponse{
		URL: url,
	}, nil
}
