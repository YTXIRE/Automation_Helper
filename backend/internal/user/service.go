package user

import (
	"backend/pkg/logging"
	"context"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *Service) Create(ctx context.Context, user CreateUserDTO) (User, error) {
	return User{}, nil
}
