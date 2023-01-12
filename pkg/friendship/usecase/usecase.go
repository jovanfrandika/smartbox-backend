package usecase

import (
	"context"

	"github.com/jovanfrandika/smartbox-backend/pkg/friendship/model"
)

func (u *usecase) CreateOne(ctx context.Context, createOneInput model.CreateOneInput) error {
	_, err := (*u.friendshipDb).CreateOne(ctx, createOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) DeleteOne(ctx context.Context, deleteOneInput model.DeleteOneInput) error {
	err := (*u.friendshipDb).DeleteOne(ctx, deleteOneInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *usecase) GetAll(ctx context.Context, userID string) (model.GetAllResponse, error) {
	friendships, err := (*u.friendshipDb).GetAll(ctx, userID)
	if err != nil {
		return model.GetAllResponse{}, err
	}

	friendIDs := []string{}
	for _, friendship := range friendships {
		if friendship.UserID != userID {
			friendIDs = append(friendIDs, friendship.UserID)
		}
		if friendship.FriendUserID != userID {
			friendIDs = append(friendIDs, friendship.FriendUserID)
		}
	}

	users, err := (*u.userDb).GetMany(ctx, friendIDs)
	if err != nil {
		return model.GetAllResponse{}, err
	}

	friends := []model.Friend{}

	for _, user := range users {
		friends = append(friends, model.Friend{
			FriendUserID: user.ID,
			Name:         user.Name,
			Email:        user.Email,
		})
	}

	return model.GetAllResponse{
		Friends: friends,
	}, nil
}
