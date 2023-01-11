package model

type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type CreateOneInput struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type GetAllResponse struct {
	Devices []Device `json:"devices"`
}
