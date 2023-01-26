package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	parcelCol "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

const (
	CODE_CHARS     = "1234567890"
	PickUpDir      = "pick-up"
	ArrivedDir     = "arrived"
	PhotoUrlFormat = "%s/%s/%s/%s"
)

func generateCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(CODE_CHARS)
	for i := 0; i < length; i++ {
		buffer[i] = CODE_CHARS[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func (u *usecase) buildFullParcel(ctx context.Context, parcels []model.Parcel) ([]model.FullParcel, error) {
	if len(parcels) <= 0 {
		return []model.FullParcel{}, nil
	}

	userIdMap := map[string]bool{}
	deviceIdMap := map[string]bool{}
	for _, parcel := range parcels {
		if _, ok := userIdMap[parcel.SenderID]; !ok {
			userIdMap[parcel.SenderID] = true
		}
		if _, ok := userIdMap[parcel.ReceiverID]; !ok && parcel.ReceiverID != parcelCol.EmptyObjectId {
			userIdMap[parcel.ReceiverID] = true
		}
		if _, ok := userIdMap[parcel.CourierID]; !ok && parcel.CourierID != parcelCol.EmptyObjectId {
			userIdMap[parcel.CourierID] = true
		}
		if _, ok := deviceIdMap[parcel.DeviceID]; !ok && parcel.DeviceID != parcelCol.EmptyObjectId {
			deviceIdMap[parcel.DeviceID] = true
		}
	}

	var userIds []string
	for userId, _ := range userIdMap {
		userIds = append(userIds, userId)
	}
	users, err := (*u.userDb).GetMany(ctx, userIds)
	if err != nil {
		return []model.FullParcel{}, err
	}

	if len(users) <= 0 {
		return []model.FullParcel{}, errors.New("empty users")
	}

	var deviceIds []string
	for deviceId, _ := range deviceIdMap {
		deviceIds = append(deviceIds, deviceId)
	}
	devices, err := (*u.deviceDb).GetMany(ctx, deviceIds)
	if err != nil {
		return []model.FullParcel{}, err
	}

	userMap := map[string]userModel.User{}
	for _, user := range users {
		if _, ok := userMap[user.ID]; !ok {
			userMap[user.ID] = user
		}
	}

	deviceMap := map[string]deviceModel.Device{}
	for _, device := range devices {
		if _, ok := deviceMap[device.ID]; !ok {
			deviceMap[device.ID] = device
		}
	}

	output := make([]model.FullParcel, len(parcels))
	for i, rawParcel := range parcels {
		parcel := model.FullParcel{
			ID:          rawParcel.ID,
			Name:        rawParcel.Name,
			Description: rawParcel.Name,
			PickUpCoor:  rawParcel.PickUpCoor,
			ArrivedCoor: rawParcel.ArrivedCoor,
			TempThr:     rawParcel.TempThr,
			HmdThr:      rawParcel.HmdThr,
			Status:      rawParcel.Status,
		}
		if rawParcel.PickUpPhoto != nil {
			parcel.PickUpPhoto = &model.FullPhoto{
				Uri:       u.buildPhotoUrl(rawParcel, model.PICK_UP_STATUS),
				UpdatedAt: rawParcel.PickUpPhoto.UpdatedAt,
			}
		}
		if rawParcel.ArrivedPhoto != nil {
			parcel.ArrivedPhoto = &model.FullPhoto{
				Uri:       u.buildPhotoUrl(rawParcel, model.ARRIVED_STATUS),
				UpdatedAt: rawParcel.ArrivedPhoto.UpdatedAt,
			}
		}
		if rawParcel.SenderID != parcelCol.EmptyObjectId {
			user, _ := userMap[rawParcel.SenderID]
			parcel.Sender = &user
		}
		if rawParcel.ReceiverID != parcelCol.EmptyObjectId {
			user, ok := userMap[rawParcel.ReceiverID]
			if ok {
				parcel.Receiver = &user
			}
		}
		if rawParcel.CourierID != parcelCol.EmptyObjectId {
			user, ok := userMap[rawParcel.CourierID]
			if ok {
				parcel.Courier = &user
			}
		}
		if rawParcel.DeviceID != parcelCol.EmptyObjectId {
			device, ok := deviceMap[rawParcel.DeviceID]
			if ok {
				parcel.Device = &device
			}
		}

		output[i] = parcel
	}

	return output, nil
}

func (u *usecase) buildPhotoUri(parcel model.Parcel, status int) string {
	switch status {
	case model.PICK_UP_STATUS:
		return fmt.Sprintf(PhotoUrlFormat, u.config.DBName, PickUpDir, parcel.SenderID, parcel.ID)
	case model.ARRIVED_STATUS:
		return fmt.Sprintf(PhotoUrlFormat, u.config.DBName, ArrivedDir, parcel.SenderID, parcel.ID)
	default:
		return ""
	}
}

func (u *usecase) buildPhotoUrl(parcel model.Parcel, status int) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.config.BucketName, u.buildPhotoUri(parcel, status))
}
