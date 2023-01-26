package usecase

import (
	"context"

	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	parcelModel "github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcelTravel/model"
)

func (u *usecase) ConsumeTravelLog(ctx context.Context, consumeTravelLogMessage model.ConsumeTravelLogMessage) error {
	device, err := (*u.deviceDb).GetOneByName(ctx, deviceModel.GetOneByNameInput{Name: consumeTravelLogMessage.DeviceName})
	if err != nil {
		return err
	}

	parcel, err := (*u.parcelDb).GetOneByDeviceAndStatus(ctx, parcelModel.GetOneByDeviceAndStatusInput{Device: device.ID, Status: parcelModel.ON_GOING_STATUS})
	if err != nil {
		return err
	}

	err = (*u.parcelTravelDb).CreateOne(ctx, model.CreateOneInput{
		ParcelID:   parcel.ID,
		Coor:       consumeTravelLogMessage.Coor,
		Temp:       consumeTravelLogMessage.Temp,
		Hmd:        consumeTravelLogMessage.Hmd,
		Sgnl:       consumeTravelLogMessage.Sgnl,
		Spd:        consumeTravelLogMessage.Spd,
		Stls:       consumeTravelLogMessage.Stls,
		DoorStatus: consumeTravelLogMessage.DoorStatus,
		GPSTs:      consumeTravelLogMessage.GPSTs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) GetAll(ctx context.Context, getAllInput model.GetAllInput) (model.GetAllResponse, error) {
	parcelTravels, err := (*u.parcelTravelDb).GetAll(ctx, getAllInput)
	if err != nil {
		return model.GetAllResponse{}, err
	}

	return model.GetAllResponse{
		ParcelTravels: parcelTravels,
	}, nil
}
