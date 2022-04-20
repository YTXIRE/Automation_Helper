package middlewares

import (
	"backend/internal/apperror"
	"backend/internal/config"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type ContextKey string

const (
	ContextUserKey ContextKey = "UserId"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		if authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apperror.NewAppError(nil, "user unauthorized", "", "UNAUTHORIZED").Marshal())
			return
		}
		headerParts := strings.Split(authorization, " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apperror.NewAppError(nil, "invalid auth header", "", "INVALID_AUTH_HEADER").Marshal())
			return
		}
		cfg := config.GetConfig()
		token, err := jwt.ParseWithClaims(headerParts[1], &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(cfg.Jwt.SecretKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apperror.NewAppError(nil, "invalid auth header", "", "INVALID_AUTH_HEADER").Marshal())
			return
		}
		claims, ok := token.Claims.(*tokenClaims)
		if !ok {
			errors.New("token claims are not of type *tokenClaims")
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, claims.UserId)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
