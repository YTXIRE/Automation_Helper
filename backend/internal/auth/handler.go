package auth

import (
	"backend/internal/apperror"
	"backend/internal/handlers"
	"backend/internal/middlewares"
	"backend/internal/user"
	"backend/pkg/logging"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

type handler struct {
	logger  *logging.Logger
	service *Service
	storage user.Storage
}

func NewHandler(logger *logging.Logger, storage user.Storage) handlers.Handler {
	return &handler{
		logger:  logger,
		storage: storage,
	}
}

const (
	singInURL       = "/sing-in"
	refreshTokenURL = "/refresh-token"
)

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, singInURL, apperror.Middleware(h.SingIn))
	router.HandlerFunc(http.MethodPost, refreshTokenURL, middlewares.AuthMiddleware(apperror.Middleware(h.RefreshToken)))
}

func (h *handler) SingIn(w http.ResponseWriter, r *http.Request) error {
	var authData DTO
	err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		return apperror.NewAppError(err, "failed to decode request body in json", "", "DECODE_ERROR")
	}
	singIn, err := h.service.SingIn(r.Context(), authData, h.storage)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(singIn)
	return nil
}

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	tokens, err := h.service.Refresh()
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "REFRESH_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
	return nil
}
