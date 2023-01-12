package model

type Coordinate struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

type Travel struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PhotoUri    string     `json:"photo_uri"`
	Start       Coordinate `json:"start"`
	End         Coordinate `json:"end"`
	ReceiverID  string     `json:"receiver_id"`
	SenderID    string     `json:"sender_id"`
	Status      int        `json:"status"`
}

type DeleteOneInput struct {
	ID string `json:"id"`
}

type CreateOneInput struct {
	SenderID string `json:"sender_id"`
}

type UpdateOneInput struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PhotoUri    string     `json:"photo_uri"`
	Start       Coordinate `json:"start"`
	End         Coordinate `json:"end"`
	ReceiverID  string     `json:"receiver_id"`
	SenderID    string     `json:"sender_id"`
	Status      int        `json:"status"`
}

type GetOneInput struct {
	ID string `json:"id"`
}

type GetPhotoSignedUrlInput struct {
	ID string `json:"id"`
}

type HistoryInput struct {
	UserID string `json:"user_id"`
}

type GetOneResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PhotoUri    string     `json:"photo_uri"`
	Start       Coordinate `json:"start"`
	End         Coordinate `json:"end"`
	ReceiverID  string     `json:"receiver_id"`
	SenderID    string     `json:"sender_id"`
	Status      int        `json:"status"`
}

type GetPhotoSignedUrlResponse struct {
	URL string `json:"url"`
}

type HistoryResponse struct {
	Histories []Travel `json:"histories"`
}
