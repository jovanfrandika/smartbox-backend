package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	parcelCol "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
	log "github.com/sirupsen/logrus"
)

func (u *usecase) GetOne(ctx context.Context, getOneInput model.GetOneInput) (model.GetOneResponse, error) {
	res, err := (*u.parcelDb).GetOne(ctx, getOneInput.ID)
	if err != nil {
		return model.GetOneResponse{}, err
	}

	fullParcels, err := u.buildFullParcel(ctx, []model.Parcel{res})
	if err != nil {
		return model.GetOneResponse{}, err
	}

	if len(fullParcels) <= 0 {
		return model.GetOneResponse{}, errors.New("Parcel Not found")
	}

	return model.GetOneResponse(fullParcels[0]), nil
}

func (u *usecase) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	_, err := (*u.parcelDb).CreateOne(ctx, createOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) UpdateOne(ctx context.Context, updateOneInput model.UpdateOneInput) error {
	err := (*u.parcelDb).UpdateOne(ctx, updateOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error {
	res, err := (*u.parcelDb).GetOne(ctx, deleteOneInput.ID)
	if err != nil {
		return err
	}

	if res.Status != parcelCol.DRAFT_STATUS {
		return errors.New("Cannot delete non draft parcel")
	}

	err = (*u.parcelDb).DeleteOne(ctx, deleteOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) Histories(ctx context.Context, historyInput model.HistoryInput) (model.HistoryResponse, error) {
	histories, err := (*u.parcelDb).Histories(ctx, historyInput)
	if err != nil {
		log.Error(fmt.Sprintf("Error: %s", err.Error()))
		return model.HistoryResponse{}, err
	}

	fullParcels, err := u.buildFullParcel(ctx, histories)
	if err != nil {
		return model.HistoryResponse{}, err
	}

	return model.HistoryResponse{
		Histories: fullParcels,
	}, nil
}

func (u *usecase) GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error) {
	res, err := (*u.parcelDb).GetOne(ctx, getPhotoSignedUrlInput.ID)
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

func (u *usecase) UpdateProgress(ctx context.Context, updateProgressInput model.UpdateProgressInput) (model.UpdateProgressResponse, error) {
	parcel, err := (*u.parcelDb).GetOne(ctx, updateProgressInput.ID)
	if err != nil {
		return model.UpdateProgressResponse{}, err
	}

	updateOneInput := model.UpdateOneInput(parcel)
	log.Debug(fmt.Sprintf("Progress %v: %v", updateOneInput.ID, updateOneInput.Status))
	switch updateOneInput.Status {
	case parcelCol.DRAFT_STATUS:
		if updateProgressInput.UserID != updateOneInput.SenderID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if len(updateOneInput.Name) <= 3 || len(updateOneInput.Description) <= 3 {
			return model.UpdateProgressResponse{}, errors.New("name or description is too short")
		}

		obj := u.storageClient.Bucket(u.config.BucketName).Object(updateOneInput.PhotoUri)
		_, err := obj.Attrs(ctx)
		if err != nil {
			return model.UpdateProgressResponse{}, fmt.Errorf("Photo not found")
		}

		if updateOneInput.Start == nil {
			return model.UpdateProgressResponse{}, errors.New("start coor can't be empty")
		}
		if updateOneInput.End == nil {
			return model.UpdateProgressResponse{}, errors.New("end coor can't be empty")
		}

		if updateOneInput.ReceiverID == parcelCol.EmptyObjectId {
			return model.UpdateProgressResponse{}, errors.New("recipient can't be empty")
		}

		updateOneInput.Status = parcelCol.WAITING_FOR_COURIER_STATUS

		err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case parcelCol.WAITING_FOR_COURIER_STATUS:
		if updateProgressInput.UserID != updateOneInput.CourierID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if updateOneInput.DeviceID == parcelCol.EmptyObjectId {
			return model.UpdateProgressResponse{}, errors.New("device can't be empty")
		}

		updateOneInput.Status = parcelCol.PICK_UP_STATUS

		err := (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case parcelCol.PICK_UP_STATUS:
		if updateProgressInput.UserID != updateOneInput.SenderID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
		err = (*u.deviceMq).PubStartTravel(device.Name, parcel.End.Lat, parcel.End.Long)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		updateOneInput.Status = parcelCol.ON_GOING_STATUS

		err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		err = (*u.deviceMq).PubStatus(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case parcelCol.ON_GOING_STATUS:
		if updateProgressInput.UserID != updateOneInput.CourierID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		// TODO: implement email verification
		updateOneInput.Status = parcelCol.ARRIVED_STATUS

		err := (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case parcelCol.ARRIVED_STATUS:
		if updateProgressInput.UserID != updateOneInput.ReceiverID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
		err = (*u.deviceMq).PubEndTravel(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		updateOneInput.Status = parcelCol.DONE_STATUS

		err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		err = (*u.deviceMq).PubStatus(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	default:
		return model.UpdateProgressResponse{}, errors.New(fmt.Sprintf("unknown status at id: %s", parcel.ID))
	}

	userIDs := []string{parcel.SenderID}
	if parcel.SenderID != parcelCol.EmptyObjectId {
		userIDs = append(userIDs, parcel.SenderID)
	}
	if parcel.CourierID != parcelCol.EmptyObjectId {
		userIDs = append(userIDs, parcel.CourierID)
	}

	users, err := (*u.userDb).GetMany(ctx, userIDs)
	if err != nil {
		return model.UpdateProgressResponse{}, err
	}

	userMap := map[string]userModel.User{}
	for _, user := range users {
		if _, ok := userMap[user.ID]; !ok {
			userMap[user.ID] = user
		}
	}

	res, err := u.buildFullParcel(ctx, []model.Parcel{model.Parcel(updateOneInput)})
	if err != nil {
		return model.UpdateProgressResponse{}, err
	}

	if len(res) <= 0 {
		return model.UpdateProgressResponse{}, errors.New("empty parcel")
	}

	return model.UpdateProgressResponse(res[0]), nil
}

func (u *usecase) OpenDoor(ctx context.Context, openDoorInput model.OpenDoorInput) error {
	parcel, err := (*u.parcelDb).GetOne(ctx, openDoorInput.ID)
	if err != nil {
		return err
	}

	userIDs := []string{parcel.SenderID}
	if parcel.SenderID != parcelCol.EmptyObjectId {
		userIDs = append(userIDs, parcel.SenderID)
	}
	if parcel.CourierID != parcelCol.EmptyObjectId {
		userIDs = append(userIDs, parcel.CourierID)
	}

	users, err := (*u.userDb).GetMany(ctx, userIDs)
	if err != nil {
		return err
	}

	userMap := map[string]userModel.User{}
	for _, user := range users {
		if _, ok := userMap[user.ID]; !ok {
			userMap[user.ID] = user
		}
	}

	switch parcel.Status {
	case parcelCol.PICK_UP_STATUS:
		if openDoorInput.UserID != parcel.SenderID {
			return errors.New("insufficient permission")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return err
		}
		err = (*u.deviceMq).PubOpenDoor(device.Name)
		if err != nil {
			return err
		}
	case parcelCol.ARRIVED_STATUS:
		if openDoorInput.UserID != parcel.ReceiverID {
			return errors.New("insufficient permission")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return err
		}
		err = (*u.deviceMq).PubOpenDoor(device.Name)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("invalid request: %s", parcel.ID))
	}

	return nil
}
