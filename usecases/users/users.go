package users

import (
	"errors"
	model "gateway-service/models"
	"gateway-service/repository/users"
	"gateway-service/util/middleware"
	"time"

	"github.com/google/uuid"
)

type svc struct {
	userStore users.UserRepository
}

func NewUserSvc(userStore users.UserRepository) *svc {
	return &svc{
		userStore: userStore,
	}
}

type UserSvc interface {
	UserRegister(req model.Users) (*uuid.UUID, error)
	UserLogin(req model.UserLoginRequest) (*model.UserLogin, error)
}

func (s *svc) UserRegister(req model.Users) (*uuid.UUID, error) {
	user, err := s.userStore.GetUserDetail(req)
	if err != nil {
		return nil, err
	}

	if user.Email == req.Email && user.Username == req.Username {
		return nil, errors.Join(errors.New("user already exists"))
	}

	salt, err := middleware.GenerateSalt(16)
	if err != nil {
		return nil, err
	}

	isPassword, err := middleware.HashPassword(req.Password, salt)
	if err != nil {
		return nil, err
	}

	req.Password = isPassword

	userID, err := s.userStore.UserRegister(req)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

func (s *svc) UserLogin(req model.UserLoginRequest) (*model.UserLogin, error) {
	user, err := s.userStore.GetUserDetail(model.Users{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	if user.Username != req.Username {
		return nil, errors.Join(errors.New("user not found"))
	}

	verifyPassword, err := middleware.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !verifyPassword {
		return nil, errors.Join(errors.New("password not match"))
	}

	tokenExpiry := time.Minute * 20
	accessToken, payload, err := middleware.CreateAccessToken(user.Email, user.Id.String(), user.Role, tokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := time.Hour * 72
	refreshToken, refreshTokenPayload, err := middleware.CreateRefreshToken(user.Email, user.Id.String(), user.Role, refreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &model.UserLogin{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiresAt.Time,
		RefreshToken:         refreshToken,
		RefreshTokenExpiryAt: refreshTokenPayload.ExpiresAt.Time,
		Users: &model.Users{
			Email:               user.Email,
			Username:            user.Username,
			Role:                user.Role,
			CategoryPreferences: user.CategoryPreferences,
			CreatedAt:           user.CreatedAt,
		},
	}, nil
}
