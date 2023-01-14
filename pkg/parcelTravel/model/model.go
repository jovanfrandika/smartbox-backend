package model

import "time"

type Coordinate struct {
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Speed      int     `json:"speed"`
	Satellites int     `json:"satellites"`
}

type ParcelTravel struct {
	ID           string     `json:"id"`
	ParcelID     string     `json:"parcel_id"`
	Coordinate   Coordinate `json:"coordinate"`
	IsDoorOpen   bool       `json:"is_door_open"`
	Signal       int        `json:"signal"`
	GPSTimestamp time.Time  `json:"gps_timestamp"`
	Timestamp    time.Time  `json:"timestamp"`
}

type CreateOneInput struct {
	ParcelID     string     `json:"parcel_id"`
	Coordinate   Coordinate `json:"coordinate"`
	IsDoorOpen   bool       `json:"is_door_open"`
	Signal       int        `json:"signal"`
	GPSTimestamp string     `json:"gps_timestamp"`
}

type GetAllInput struct {
	ParcelID string `json:"parcel_id"`
}

type GetAllResponse struct {
	ParcelTravels []ParcelTravel `json:"parcel_travels"`
}

type ConsumeTravelLogMessage struct {
	DeviceName   string     `json:"device_name"`
	Coordinate   Coordinate `json:"coordinate"`
	IsDoorOpen   bool       `json:"is_door_open"`
	Signal       int        `json:"signal"`
	GPSTimestamp string     `json:"gps_timestamp"`
}
