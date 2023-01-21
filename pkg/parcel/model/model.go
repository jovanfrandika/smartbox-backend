package model

import (
	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

const (
	DRAFT_STATUS               = int(0)
	WAITING_FOR_COURIER_STATUS = int(1)
	PICK_UP_STATUS             = int(2)
	ON_GOING_STATUS            = int(3)
	ARRIVED_STATUS             = int(4)
	DONE_STATUS                = int(5)
)

type Coordinate struct {
	Lat   float32 `json:"lat"`
	Long  float32 `json:"long"`
	Temp  float32 `json:"temp"`
	Humid float32 `json:"humid"`
}

type Parcel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	PhotoUri     string      `json:"photo_uri"`
	IsPhotoValid bool        `json:"is_photo_valid"`
	Start        *Coordinate `json:"start"`
	End          *Coordinate `json:"end"`
	ReceiverID   string      `json:"receiver_id"`
	SenderID     string      `json:"sender_id"`
	CourierID    string      `json:"courier_id"`
	DeviceID     string      `json:"device_id"`
	Status       int         `json:"status"`
}

type FullParcel struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	PhotoUri     string              `json:"photo_uri"`
	IsPhotoValid bool                `json:"is_photo_valid"`
	Start        *Coordinate         `json:"start"`
	End          *Coordinate         `json:"end"`
	Receiver     *userModel.User     `json:"receiver"`
	Sender       *userModel.User     `json:"sender"`
	Courier      *userModel.User     `json:"courier"`
	Device       *deviceModel.Device `json:"device"`
	Status       int                 `json:"status"`
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

type GetOneByDeviceAndStatusInput struct {
	Device string `json:"device"`
	Status int    `json:"status"`
}

type GetPhotoSignedUrlInput struct {
	ID string `json:"id"`
}

type CheckPhotoInput struct {
	ID string `json:"id"`
}

type SendParcelCodeToReceiverInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type VerifyParcelCodeInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Code   string `json:"code"`
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
