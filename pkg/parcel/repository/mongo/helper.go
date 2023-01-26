package rMongo

import "github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"

func buildParcel(parcel Parcel) model.Parcel {
	var pickUpCoor *model.Coordinate
	if parcel.PickUpLoc != nil {
		pickUpCoor = &model.Coordinate{
			Lat: parcel.PickUpLoc.Coordinates[1],
			Lng: parcel.PickUpLoc.Coordinates[0],
		}
	}
	var arrivedCoor *model.Coordinate
	if parcel.ArrivedLoc != nil {
		arrivedCoor = &model.Coordinate{
			Lat: parcel.ArrivedLoc.Coordinates[1],
			Lng: parcel.ArrivedLoc.Coordinates[0],
		}
	}
	var pickUpPhoto *model.Photo
	if parcel.PickUpPhoto != nil {
		pickUpPhoto = &model.Photo{
			UpdatedAt: parcel.PickUpPhoto.UpdatedAt.Time(),
		}
	}
	var arrivedPhoto *model.Photo
	if parcel.ArrivedPhoto != nil {
		arrivedPhoto = &model.Photo{
			UpdatedAt: parcel.ArrivedPhoto.UpdatedAt.Time(),
		}
	}
	var tempThr *model.Threshold
	if parcel.TempThr != nil {
		tempThr = &model.Threshold{
			Low:  parcel.TempThr.High,
			High: parcel.TempThr.Low,
		}
	}
	var hmdThr *model.Threshold
	if parcel.HmdThr != nil {
		hmdThr = &model.Threshold{
			Low:  parcel.HmdThr.High,
			High: parcel.HmdThr.Low,
		}
	}

	return model.Parcel{
		ID:           parcel.ID.Hex(),
		Name:         parcel.Name,
		Description:  parcel.Description,
		PickUpCoor:   pickUpCoor,
		ArrivedCoor:  arrivedCoor,
		PickUpPhoto:  pickUpPhoto,
		ArrivedPhoto: arrivedPhoto,
		TempThr:      tempThr,
		HmdThr:       hmdThr,
		ReceiverID:   parcel.ReceiverID.Hex(),
		SenderID:     parcel.SenderID.Hex(),
		CourierID:    parcel.CourierID.Hex(),
		DeviceID:     parcel.DeviceID.Hex(),
		Status:       parcel.Status,
	}
}
