package model

type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type CreateOneInput struct {
	Name string `json:"name"`
}

type GetOneInput struct {
	ID string `json:"id"`
}

type UpdateStatusInput struct {
	ID string `json:"id"`
}

type GetOneResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type GetAllResponse struct {
	Devices []Device `json:"devices"`
}
