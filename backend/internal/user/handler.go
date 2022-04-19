package user

import (
	"backend/internal/apperror"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/pkg/logging"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

type handler struct {
	logger  *logging.Logger
	service *Service
	storage Storage
}

func NewHandler(logger *logging.Logger, storage Storage) handlers.Handler {
	return &handler{
		logger:  logger,
		storage: storage,
	}
}

const (
	usersUrl = "/users"
	userURL  = "/users/:uuid"
)

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersUrl, middleware.AuthMiddleware(apperror.Middleware(h.GetList)))
	router.HandlerFunc(http.MethodGet, userURL, middleware.AuthMiddleware(apperror.Middleware(h.GetUserByUUID)))
	router.HandlerFunc(http.MethodPost, usersUrl, middleware.AuthMiddleware(apperror.Middleware(h.CreateUser)))
	router.HandlerFunc(http.MethodPut, userURL, middleware.AuthMiddleware(apperror.Middleware(h.UpdateUser)))
	router.HandlerFunc(http.MethodDelete, userURL, middleware.AuthMiddleware(apperror.Middleware(h.DeleteUser)))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	userList, err := h.service.GetUserList(context.Background(), h.storage)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "GET_USER_LIST_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	oid := params.ByName("uuid")
	user, err := h.service.GetUserByID(context.Background(), h.storage, oid)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "GET_FIND_BY_ID_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var newUser DTO
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return apperror.NewAppError(err, "failed to decode request body in json", "", "DECODE_ERROR")
	}
	user, err := h.service.Create(context.Background(), newUser, h.storage)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "CREATE_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	var userData DTO
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return apperror.NewAppError(err, "failed to decode request body in json", "", "DECODE_ERROR")
	}
	params := httprouter.ParamsFromContext(r.Context())
	userData.ID = params.ByName("uuid")
	user, err := h.service.UpdateUser(context.Background(), h.storage, userData)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "UPDATE_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	oid := params.ByName("uuid")
	err := h.service.DeleteUser(context.Background(), h.storage, oid)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "DELETE_ERROR")
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
