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

type GetOneByNameInput struct {
	Name string `json:"name"`
}

type ConsumeUpdateStatusMessage struct {
	Name        string `json:"name"`
	Status      int    `json:"status"`
	LogInterval int    `json:"log_interval"`
}

type UpdateStatusInput struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type GetOneResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type GetAllResponse struct {
	Devices []Device `json:"devices"`
}
