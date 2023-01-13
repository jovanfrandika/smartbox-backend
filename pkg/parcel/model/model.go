package model

import (
	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

type Coordinate struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

type Parcel struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	PhotoUri    string      `json:"photo_uri"`
	Start       *Coordinate `json:"start"`
	End         *Coordinate `json:"end"`
	ReceiverID  string      `json:"receiver_id"`
	SenderID    string      `json:"sender_id"`
	CourierID   string      `json:"courier_id"`
	DeviceID    string      `json:"device_id"`
	Status      int         `json:"status"`
}

type FullParcel struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	PhotoUri    string              `json:"photo_uri"`
	Start       *Coordinate         `json:"start"`
	End         *Coordinate         `json:"end"`
	Receiver    *userModel.User     `json:"receiver"`
	Sender      *userModel.User     `json:"sender"`
	Courier     *userModel.User     `json:"courier"`
	Device      *deviceModel.Device `json:"device"`
	Status      int                 `json:"status"`
}

type DeleteOneInput struct {
	ID string `json:"id"`
}

type CreateOneInput struct {
	SenderID string `json:"sender_id"`
}

type UpdateOneInput = Parcel

type GetOneInput struct {
	ID string `json:"id"`
}

type GetPhotoSignedUrlInput struct {
	ID string `json:"id"`
}

type UpdateProgressInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type HistoryInput struct {
	UserID string `json:"user_id"`
}

type OpenDoorInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type GetOneResponse = FullParcel

type GetPhotoSignedUrlResponse struct {
	URL string `json:"url"`
}

type UpdateProgressResponse = FullParcel

type HistoryResponse struct {
	Histories []FullParcel `json:"histories"`
}
