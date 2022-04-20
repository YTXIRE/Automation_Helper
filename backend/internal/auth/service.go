package auth

import (
	"backend/internal/apperror"
	"backend/internal/config"
	"backend/internal/user"
	"backend/pkg/encrypt"
	"backend/pkg/logging"
	"backend/pkg/times"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Service struct {
	logger *logging.Logger
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func (s *Service) SingIn(ctx context.Context, auth DTO, storage user.Storage) (Respond, error) {
	err := auth.Validate()
	if err != nil {
		return Respond{}, apperror.NewAppError(err, err.Error(), "", "VALIDATE_ERROR")
	}
	findUser, err := storage.FindByField(ctx, "username", auth.Username)
	if err != nil {
		return Respond{}, apperror.NewAppError(err, fmt.Sprintf("failed to find user with username: %s. error: %v", auth.Username, err), "", "FIND_ERROR")
	}
	if !encrypt.CheckPasswordHash(auth.Password, findUser.PasswordHash) {
		return Respond{}, apperror.NewAppError(nil, "passwords don't match", "", "ERROR_CHECK_PASSWORD")
	}
	cfg := config.GetConfig()
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(times.ConvertTime(cfg.Jwt.AccessToken.TTL, cfg.Jwt.AccessToken.Type)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		findUser.ID,
	})
	accessToken, err := accessTokenClaims.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return Respond{}, apperror.NewAppError(err, err.Error(), "", "GENERATE_ACCESS_TOKEN_ERROR")
	}
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(times.ConvertTime(cfg.Jwt.RefreshToken.TTL, cfg.Jwt.RefreshToken.Type)).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	refreshToken, err := refreshTokenClaims.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return Respond{}, apperror.NewAppError(err, err.Error(), "", "GENERATE_REFRESH_TOKEN_ERROR")
	}
	generateTokens := Tokens{
		UserId:       findUser.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	findUser.LastLogin = time.Now().Unix()
	err = storage.Update(ctx, findUser)
	if err != nil {
		return Respond{}, apperror.NewAppError(err, err.Error(), "", "USER_UPDATE_ERROR")
	}
	return Respond{
		ID:           findUser.ID,
		Username:     findUser.Username,
		AccessToken:  generateTokens.AccessToken,
		RefreshToken: generateTokens.RefreshToken,
		Email:        findUser.Email,
		CreatedAt:    findUser.CreatedAt,
		UpdatedAt:    findUser.UpdatedAt,
		LastLogin:    findUser.LastLogin,
	}, nil
}

func (s *Service) Refresh() (Tokens, error) {
	cfg := config.GetConfig()
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(times.ConvertTime(cfg.Jwt.AccessToken.TTL, cfg.Jwt.AccessToken.Type)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		"1",
	})
	accessToken, err := accessTokenClaims.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return Tokens{}, apperror.NewAppError(err, err.Error(), "", "GENERATE_ACCESS_TOKEN_ERROR")
	}
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(times.ConvertTime(cfg.Jwt.RefreshToken.TTL, cfg.Jwt.RefreshToken.Type)).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	refreshToken, err := refreshTokenClaims.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return Tokens{}, apperror.NewAppError(err, err.Error(), "", "GENERATE_REFRESH_TOKEN_ERROR")
	}
	generateTokens := Tokens{
		UserId:       "1",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return generateTokens, nil
}
