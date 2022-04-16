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

func (s *Service) Create(ctx context.Context, user DTO, storage Storage) (User, error) {
	err := user.Validate()
	if err != nil {
		return User{}, apperror.NewAppError(err, err.Error(), "", "US-000004")
	}
	findUser, err := storage.FindByField(ctx, "email", user.Email)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with email: %s. error: %v", user.Email, err), "", "US-000005")
	}
	if findUser.Email != "" {
		return User{}, apperror.NewAppError(nil, "this email is already busy", "", "US-000006")
	}
	findUser, err = storage.FindByField(ctx, "username", user.Username)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with username: %s. error: %v", user.Username, err), "", "US-000008")
	}
	if findUser.Username != "" {
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

func (s *Service) GetUserList(ctx context.Context, storage Storage) ([]User, error) {
	users, err := storage.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, err.Error(), "", "US-000013")
	}
	if len(users) == 0 {
		return []User{}, nil
	}
	return users, nil
}

func (s *Service) GetUserByID(ctx context.Context, storage Storage, id string) (User, error) {
	user, err := storage.FindOne(ctx, id)
	if err != nil {
		return User{}, apperror.ErrNotFound
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, storage Storage, user DTO) (User, error) {
	err := user.Validate()
	if err != nil {
		return User{}, apperror.NewAppError(err, err.Error(), "", "US-000004")
	}
	findUser, err := storage.FindByField(ctx, "email", user.Email)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with email: %s. error: %v", user.Email, err), "", "US-000005")
	}
	if findUser.Email != "" && findUser.ID != user.ID {
		return User{}, apperror.NewAppError(nil, "this email is already busy", "", "US-000006")
	}
	findUser, err = storage.FindByField(ctx, "username", user.Username)
	if err != nil {
		return User{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with username: %s. error: %v", user.Username, err), "", "US-000008")
	}
	if findUser.Username != "" && findUser.ID != user.ID {
		return User{}, apperror.NewAppError(nil, "this username is already busy", "", "US-000007")
	}
	findUser, err = storage.FindOne(ctx, user.ID)
	if err != nil {
		return User{}, apperror.ErrNotFound
	}
	if user.Email != "" && findUser.Email != user.Email {
		findUser.Email = user.Email
	}
	if user.Username != "" && findUser.Username != user.Username {
		findUser.Username = user.Username
	}
	if user.Password != "" {
		hash, err := encrypt.HashPassword(user.Password)
		if err != nil {
			return User{}, apperror.NewAppError(err, err.Error(), "", "US-000010")
		}
		findUser.PasswordHash = hash
	}
	findUser.UpdatedAt = time.Now().Unix()
	err = storage.Update(ctx, findUser)
	if err != nil {
		return User{}, apperror.NewAppError(err, "failed to update user", "", "US-000017")
	}
	return findUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, storage Storage, id string) error {
	err := storage.Delete(ctx, id)
	if err != nil {
		return apperror.ErrNotFound
	}
	return nil
}
