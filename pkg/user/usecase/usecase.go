package usecase

import (
	"context"
	"time"

	utils "github.com/jovanfrandika/smartbox-backend/pkg/common/jwt"
	"github.com/jovanfrandika/smartbox-backend/pkg/user/model"
)

const (
	ACCESS_TOKEN_EXPIRES_IN  = time.Duration(3600 * 24 * 30 * 1 * time.Second)
	REFRESH_TOKEN_EXPIRES_IN = time.Duration(3600 * 24 * 30 * 3 * time.Second)
)

func (u *usecase) Me(ctx context.Context, meInput model.MeInput) (model.MeResponse, error) {
	user, err := (*u.db).GetUser(ctx, meInput.ID)
	if err != nil {
		return model.MeResponse{}, err
	}

	return model.MeResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (u *usecase) Search(ctx context.Context, searchInput model.SearchInput) (model.SearchResponse, error) {
	users, err := (*u.db).Search(ctx, searchInput.Email)
	if err != nil {
		return model.SearchResponse{}, err
	}

	return model.SearchResponse{
		Users: users,
	}, nil
}

func (u *usecase) Register(ctx context.Context, registerInput model.RegisterInput) (model.RegisterResponse, error) {
	userID, err := (*u.db).CreateUser(ctx, registerInput)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	accessToken, err := utils.CreateToken(ACCESS_TOKEN_EXPIRES_IN, userID, u.config.JWTAccessSecretKey)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	refreshToken, err := utils.CreateToken(REFRESH_TOKEN_EXPIRES_IN, userID, u.config.JWTRefreshSecretKey)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	return model.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func (u *usecase) Login(ctx context.Context, loginInput model.LoginInput) (model.LoginResponse, error) {
	user, err := (*u.db).Login(ctx, loginInput)
	if err != nil {
		return model.LoginResponse{}, err
	}

	accessToken, err := utils.CreateToken(ACCESS_TOKEN_EXPIRES_IN, user.ID, u.config.JWTAccessSecretKey)
	if err != nil {
		return model.LoginResponse{}, err
	}

	refreshToken, err := utils.CreateToken(REFRESH_TOKEN_EXPIRES_IN, user.ID, u.config.JWTRefreshSecretKey)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *usecase) RefreshAccessToken(ctx context.Context, refreshInput model.RefreshInput) (model.RefreshResponse, error) {
	userId, err := utils.ValidateToken(refreshInput.RefreshToken, u.config.JWTRefreshSecretKey)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	user, err := (*u.db).GetUser(ctx, userId)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	accessToken, err := utils.CreateToken(ACCESS_TOKEN_EXPIRES_IN, user.ID, u.config.JWTAccessSecretKey)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	return model.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}
