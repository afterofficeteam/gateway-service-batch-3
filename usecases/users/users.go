package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	model "gateway-service/models"
	"gateway-service/repository/users"
	"gateway-service/util/middleware"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type svc struct {
	userStore users.UserRepository
	redis     *redis.Client
}

func NewUserSvc(userStore users.UserRepository, redis *redis.Client) *svc {
	return &svc{
		userStore: userStore,
		redis:     redis,
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
	ctx := context.Background()

	val, err := s.redis.Get(ctx, req.Username).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis: Key %s not found", req.Username)
		}
	}

	if val != "" {
		var userLogin model.UserLogin
		if err := json.Unmarshal([]byte(val), &userLogin); err != nil {
			return nil, errors.Join(fmt.Errorf("error unmarshaling Redis value for %s: %v", req.Username, err))
		}

		return &userLogin, nil
	}

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

	userLogin := model.UserLogin{
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
	}

	userData, err := json.Marshal(userLogin)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, req.Username, userData, 20*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &userLogin, nil
}
