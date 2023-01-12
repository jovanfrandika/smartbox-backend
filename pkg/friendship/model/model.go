package model

type Friendship struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	FriendUserID string `json:"friend_user_id"`
}

type Friend struct {
	FriendUserID string `json:"friend_user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}

type CreateOneInput struct {
	UserID       string `json:"user_id"`
	FriendUserID string `json:"friend_user_id"`
}

type DeleteOneInput struct {
	UserID       string `json:"user_id"`
	FriendUserID string `json:"friend_user_id"`
}

type GetAllResponse struct {
	Friends []Friend `json:"friends"`
}
