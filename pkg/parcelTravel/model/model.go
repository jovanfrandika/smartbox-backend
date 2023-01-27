package model

import "time"

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ParcelTravel struct {
	ID         string     `json:"id"`
	ParcelID   string     `json:"parcel_id"`
	Coor       Coordinate `json:"coor"`
	Temp       float32    `json:"temp"`
	Hmd        float32    `json:"hmd"`
	DoorStatus int        `json:"door_status"`
	Sgnl       int        `json:"sgnl"`
	Spd        int        `json:"spd"`
	Stls       int        `json:"stls"`
	GPSTs      time.Time  `json:"gps_ts"`
	Ts         time.Time  `json:"ts"`
}

type CreateOneInput struct {
	ParcelID   string     `json:"parcel_id"`
	Coor       Coordinate `json:"coor"`
	Temp       float32    `json:"temp"`
	Hmd        float32    `json:"hmd"`
	DoorStatus int        `json:"door_status"`
	Sgnl       int        `json:"sgnl"`
	Spd        int        `json:"spd"`
	Stls       int        `json:"stls"`
	GPSTs      string     `json:"gps_ts"`
}

type GetAllInput struct {
	ParcelID string `json:"parcel_id"`
}

type GetAllResponse struct {
	ParcelTravels []ParcelTravel `json:"parcel_travels"`
}

type ConsumeTravelLogMessage struct {
	DeviceName string     `json:"device_name"`
	Coor       Coordinate `json:"coor"`
	Temp       float32    `json:"temp"`
	Hmd        float32    `json:"hmd"`
	DoorStatus int        `json:"door_status"`
	Sgnl       int        `json:"sgnl"`
	Spd        int        `json:"spd"`
	Stls       int        `json:"stls"`
	GPSTs      string     `json:"gps_ts"`
}
