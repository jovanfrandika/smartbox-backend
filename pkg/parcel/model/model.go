package model

import (
	"time"

	deviceModel "github.com/jovanfrandika/smartbox-backend/pkg/device/model"
	userModel "github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

const (
	DRAFT_STATUS               = int(1)
	WAITING_FOR_COURIER_STATUS = int(2)
	PICK_UP_STATUS             = int(3)
	ON_GOING_STATUS            = int(4)
	ARRIVED_STATUS             = int(5)
	DONE_STATUS                = int(6)
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Photo struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type FullPhoto struct {
	Uri       string    `json:"uri"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Threshold struct {
	Low  float32 `json:"low"`
	High float32 `json:"high"`
}

type Parcel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	PickUpPhoto  *Photo      `json:"pick_up_photo"`
	ArrivedPhoto *Photo      `json:"arrived_photo"`
	PickUpCoor   *Coordinate `json:"pick_up_coor"`
	ArrivedCoor  *Coordinate `json:"arrived_coor"`
	TempThr      *Threshold  `json:"temp_thr"`
	HmdThr       *Threshold  `json:"hmd_thr"`
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
	PickUpPhoto  *FullPhoto          `json:"pick_up_photo"`
	ArrivedPhoto *FullPhoto          `json:"arrived_photo"`
	PickUpCoor   *Coordinate         `json:"pick_up_coor"`
	ArrivedCoor  *Coordinate         `json:"arrived_coor"`
	TempThr      *Threshold          `json:"temp_thr"`
	HmdThr       *Threshold          `json:"hmd_thr"`
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
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type CheckPhotoInput struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type SendParcelCodeInput struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	ToUserID string `json:"to_user_id"`
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

type GetNearbyPickUpsInput struct {
	UserCoor Coordinate `json:"user_coor"`
}

type HistoryInput struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type OpenDoorInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CloseDoorInput struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type GetOneResponse = FullParcel

type GetPhotoSignedUrlResponse struct {
	URL string `json:"url"`
}

type GetNearbyPickUpsResponse struct {
	Parcels []FullParcel `json:"parcels"`
}

type UpdateProgressResponse = FullParcel

type HistoryResponse struct {
	Histories []FullParcel `json:"histories"`
}
