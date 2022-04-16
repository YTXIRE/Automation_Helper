package user

import (
	"backend/internal/apperror"
	"backend/pkg/encrypt"
	"backend/pkg/logging"
	"context"
	"fmt"
	"time"
)

type Service struct {
	logger *logging.Logger
}

func (s *Service) Create(ctx context.Context, user CreateUserDTO, storage Storage) (User, error) {
	err := user.Validate()
	if err != nil {
		return User{}, apperror.NewAppError(err, err.Error(), "", "US-000004")
	}
	users, err := storage.FindByField(ctx, "email", user.Email)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with email: %s. error: %v", user.Email, err), "", "US-000005")
	}
	if users.Email != "" {
		return User{}, apperror.NewAppError(nil, "this email is already busy", "", "US-000006")
	}
	users, err = storage.FindByField(ctx, "username", user.Username)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with username: %s. error: %v", user.Username, err), "", "US-000008")
	}
	if users.Username != "" {
		return User{}, apperror.NewAppError(nil, "this username is already busy", "", "US-000007")
	}
	hash, err := encrypt.HashPassword(user.Password)
	if err != nil {
		return User{}, apperror.NewAppError(err, err.Error(), "", "US-000010")
	}
	createUserData := User{
		ID:           "",
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: hash,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    0,
		LastLogin:    0,
	}
	oid, err := storage.Create(ctx, createUserData)
	if err != nil {
		return User{}, apperror.NewAppError(err, err.Error(), "", "US-000011")
	}
	createUserData.ID = oid
	return createUserData, err
}
