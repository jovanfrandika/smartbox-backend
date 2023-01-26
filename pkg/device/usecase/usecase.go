package usecase

import (
	"context"
	"fmt"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/qr"
	"github.com/jovanfrandika/smartbox-backend/pkg/device/model"
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
		return model.GetAllResponse{}, err
	}

	return model.GetAllResponse{
		Devices: devices,
	}, nil
}

func (u *usecase) GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error) {
	device, err := (*u.db).GetOne(ctx, getOneInput)
	if err != nil {
		return model.GetOneResponse{}, err
	}
	return model.GetOneResponse(device), nil
}

func (u *usecase) GetQRCode(ctx context.Context, getQRCodeInput model.GetQRCodeInput) (model.GetQRCodeResponse, error) {
	device, err := (*u.db).GetOne(ctx, model.GetOneInput{
		ID: getQRCodeInput.ID,
	})
	if err != nil {
		return model.GetQRCodeResponse{}, err
	}

	encoded, err := qr.EncodeStringToPng(fmt.Sprintf("%s", device.ID))
	if err != nil {
		return model.GetQRCodeResponse{}, err
	}

	return model.GetQRCodeResponse{
		QRCode: encoded,
	}, nil
}
