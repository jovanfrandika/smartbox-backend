package usecase

import (
	"context"
	"fmt"

	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	log "github.com/sirupsen/logrus"
)

func (u *usecase) ConsumeUpdateStatus(ctx context.Context, consumeUpdateStatusInput model.ConsumeUpdateStatusMessage) error {
	err := (*u.db).UpdateStatus(ctx, model.UpdateStatusInput{
		Name:   consumeUpdateStatusInput.Name,
		Status: consumeUpdateStatusInput.Status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	_, err := (*u.db).CreateOne(ctx, createOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) GetAll(ctx context.Context) (model.GetAllResponse, error) {
	devices, err := (*u.db).GetAll(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error: %s", err.Error()))
		return model.GetAllResponse{}, err
	}

	return model.GetAllResponse{
		Devices: devices,
	}, nil
}
