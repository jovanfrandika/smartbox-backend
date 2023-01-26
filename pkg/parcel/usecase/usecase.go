package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/email"
	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	"github.com/jovanfrandika/smartbox-backend/pkg/parcel/model"
	parcelCol "github.com/jovanfrandika/smartbox-backend/pkg/parcel/repository/mongo"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

const (
	CODE_EMAIL_SUBJECT = "Kode membuka untuk paket %v"
	CODE_EMAIL_BODY    = "Berikut adalah kode untuk membuka paket %v\n\t %v \n Terima kasih,\nSalam admin"
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

	if res.Status != model.DRAFT_STATUS {
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

func (u *usecase) GetNearbyPickUps(ctx context.Context, getNearyPickUpsInput model.GetNearbyPickUpsInput) (model.GetNearbyPickUpsResponse, error) {
	pickUps, err := (*u.parcelDb).GetNearbyPickUps(ctx, getNearyPickUpsInput)
	if err != nil {
		return model.GetNearbyPickUpsResponse{}, err
	}

	fullParcels, err := u.buildFullParcel(ctx, pickUps)
	if err != nil {
		return model.GetNearbyPickUpsResponse{}, err
	}

	return model.GetNearbyPickUpsResponse{
		Parcels: fullParcels,
	}, nil
}

func (u *usecase) GetPhotoSignedUrl(ctx context.Context, getPhotoSignedUrlInput model.GetPhotoSignedUrlInput) (model.GetPhotoSignedUrlResponse, error) {
	res, err := (*u.parcelDb).GetOne(ctx, getPhotoSignedUrlInput.ID)
	if err != nil {
		return model.GetPhotoSignedUrlResponse{}, err
	}

	url, err := (*u.storage).GetSignedUrl(u.buildPhotoUri(res, getPhotoSignedUrlInput.Status))
	if err != nil {
		return model.GetPhotoSignedUrlResponse{}, err
	}

	return model.GetPhotoSignedUrlResponse{
		URL: url,
	}, nil
}

func (u *usecase) CheckPhoto(ctx context.Context, checkPhotoInput model.CheckPhotoInput) error {
	res, err := (*u.parcelDb).GetOne(ctx, checkPhotoInput.ID)
	if err != nil {
		return err
	}

	updateOneInput := model.UpdateOneInput(res)
	updatedAt, err := (*u.storage).GetObjectUpdatedAt(ctx, u.buildPhotoUri(res, checkPhotoInput.Status))
	if err != nil {
		return err
	}

	if checkPhotoInput.Status == model.PICK_UP_STATUS {
		updateOneInput.PickUpPhoto = &model.Photo{
			UpdatedAt: updatedAt,
		}
	} else if checkPhotoInput.Status == model.ARRIVED_STATUS {
		updateOneInput.ArrivedPhoto = &model.Photo{
			UpdatedAt: updatedAt,
		}
	}

	err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) SendParcelCode(ctx context.Context, sendParcelCodeInput model.SendParcelCodeInput) error {
	parcel, err := (*u.parcelDb).GetOne(ctx, sendParcelCodeInput.ID)
	if err != nil {
		return err
	}

	if sendParcelCodeInput.UserID != parcel.CourierID {
		return errors.New("Insufficient Permission")
	}

	if !(sendParcelCodeInput.UserID == parcel.ReceiverID || sendParcelCodeInput.UserID == parcel.SenderID) {
		return errors.New("Wrong code recipient")
	}

	code, _ := (*u.parcelCache).GetParcelCode(ctx, parcel.ID)
	if code != "" {
		return errors.New("Parcel code already exists")
	}

	code, err = generateCode(6)
	if err != nil {
		return errors.New("Generate code failed")
	}

	err = (*u.parcelCache).SetParcelCode(ctx, parcel.ID, code)
	if err != nil {
		return errors.New("Generate code failed")
	}

	users, err := (*u.userDb).GetMany(ctx, []string{sendParcelCodeInput.ToUserID})
	if err != nil || len(users) <= 0 {
		return errors.New("User not found")
	}

	err = (*u.email).Send(ctx, email.SendInput{
		To:      users[0].Email,
		Subject: fmt.Sprintf(CODE_EMAIL_SUBJECT, parcel.Name),
		Body:    fmt.Sprintf(CODE_EMAIL_BODY, parcel.Name, code),
	})
	if err != nil {
		return errors.New("Send code failed")
	}

	return nil
}

func (u *usecase) VerifyParcelCode(ctx context.Context, verifyParcelCodeInput model.VerifyParcelCodeInput) error {
	parcel, err := (*u.parcelDb).GetOne(ctx, verifyParcelCodeInput.ID)
	if err != nil {
		return err
	}

	if !(verifyParcelCodeInput.UserID == parcel.ReceiverID || verifyParcelCodeInput.UserID == parcel.SenderID) {
		return errors.New("Insufficient permission")
	}

	code, err := (*u.parcelCache).GetParcelCode(ctx, parcel.ID)
	if err != nil {
		return errors.New("Code already expired")
	}

	if verifyParcelCodeInput.Code != code {
		return errors.New("Invalid code")
	}

	updateOneInput := model.UpdateOneInput(parcel)
	updateOneInput.Status = model.ARRIVED_STATUS
	err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) UpdateProgress(ctx context.Context, updateProgressInput model.UpdateProgressInput) (model.UpdateProgressResponse, error) {
	parcel, err := (*u.parcelDb).GetOne(ctx, updateProgressInput.ID)
	if err != nil {
		return model.UpdateProgressResponse{}, err
	}

	updateOneInput := model.UpdateOneInput(parcel)
	switch updateOneInput.Status {
	case model.DRAFT_STATUS:
		if updateProgressInput.UserID != updateOneInput.SenderID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if len(updateOneInput.Name) <= 3 {
			return model.UpdateProgressResponse{}, errors.New("Name is too short")
		}

		if len(updateOneInput.Description) <= 3 {
			return model.UpdateProgressResponse{}, errors.New("Description is too short")
		}

		if updateOneInput.PickUpCoor == nil {
			return model.UpdateProgressResponse{}, errors.New("Pick up coor can't be empty")
		}
		if updateOneInput.ArrivedCoor == nil {
			return model.UpdateProgressResponse{}, errors.New("Arrived coor can't be empty")
		}

		if updateOneInput.ReceiverID == parcelCol.EmptyObjectId {
			return model.UpdateProgressResponse{}, errors.New("recipient can't be empty")
		}

		updateOneInput.Status = model.WAITING_FOR_COURIER_STATUS

		err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case model.WAITING_FOR_COURIER_STATUS:
		if updateProgressInput.UserID != updateOneInput.CourierID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if updateOneInput.DeviceID == parcelCol.EmptyObjectId {
			return model.UpdateProgressResponse{}, errors.New("Device can't be empty")
		}

		updateOneInput.Status = model.PICK_UP_STATUS

		err := (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case model.PICK_UP_STATUS:
		if updateProgressInput.UserID != updateOneInput.CourierID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if updateOneInput.PickUpPhoto == nil {
			return model.UpdateProgressResponse{}, errors.New("Pick up Photo can't be empty")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
		err = (*u.deviceMq).PubStartTravel(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		updateOneInput.Status = model.ON_GOING_STATUS

		err = (*u.parcelDb).UpdateOne(ctx, updateOneInput)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		err = (*u.deviceMq).PubStatus(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
	case model.ON_GOING_STATUS:
		return model.UpdateProgressResponse{}, errors.New("Not implemented")
	case model.ARRIVED_STATUS:
		if updateProgressInput.UserID != updateOneInput.ReceiverID {
			return model.UpdateProgressResponse{}, errors.New("insufficient permission")
		}

		if updateOneInput.ArrivedPhoto == nil {
			return model.UpdateProgressResponse{}, errors.New("Arrived Photo can't be empty")
		}

		device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}
		err = (*u.deviceMq).PubEndTravel(device.Name)
		if err != nil {
			return model.UpdateProgressResponse{}, err
		}

		updateOneInput.Status = model.DONE_STATUS

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
	case model.PICK_UP_STATUS:
		if openDoorInput.UserID != parcel.SenderID {
			return errors.New("insufficient permission")
		}
	case model.ARRIVED_STATUS:
		if openDoorInput.UserID != parcel.ReceiverID {
			return errors.New("insufficient permission")
		}
	default:
		return errors.New(fmt.Sprintf("invalid request: %s", parcel.ID))
	}

	device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
	if err != nil {
		return err
	}
	err = (*u.deviceMq).PubOpenDoor(device.Name)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) CloseDoor(ctx context.Context, closeDoorInput model.CloseDoorInput) error {
	parcel, err := (*u.parcelDb).GetOne(ctx, closeDoorInput.ID)
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
	case model.PICK_UP_STATUS:
		if closeDoorInput.UserID != parcel.SenderID {
			return errors.New("insufficient permission")
		}
	case model.ARRIVED_STATUS:
		if closeDoorInput.UserID != parcel.ReceiverID {
			return errors.New("insufficient permission")
		}
	default:
		return errors.New(fmt.Sprintf("invalid request: %s", parcel.ID))
	}

	device, err := (*u.deviceDb).GetOne(ctx, deviceModel.GetOneInput{ID: parcel.DeviceID})
	if err != nil {
		return err
	}
	err = (*u.deviceMq).PubCloseDoor(device.Name)
	if err != nil {
		return err
	}

	return nil
}
