package project

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
	projectsUrl = "/projects"
	projectURL  = "/projects/:uuid"
)

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, projectsUrl, middleware.AuthMiddleware(apperror.Middleware(h.GetList)))
	router.HandlerFunc(http.MethodGet, projectURL, middleware.AuthMiddleware(apperror.Middleware(h.GetProjectByUUID)))
	router.HandlerFunc(http.MethodPost, projectsUrl, middleware.AuthMiddleware(apperror.Middleware(h.CreateProject)))
	router.HandlerFunc(http.MethodPut, projectURL, middleware.AuthMiddleware(apperror.Middleware(h.UpdateProject)))
	router.HandlerFunc(http.MethodDelete, projectURL, middleware.AuthMiddleware(apperror.Middleware(h.DeleteProject)))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	projectsList, err := h.service.GetProjectsList(context.Background(), h.storage)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "GET_PROJECT_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectsList)
	return nil
}

func (h *handler) GetProjectByUUID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	oid := params.ByName("uuid")
	user, err := h.service.GetProjectByID(context.Background(), h.storage, oid)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "GET_PROJECT_BY_ID_ERROR")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) CreateProject(w http.ResponseWriter, r *http.Request) error {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		return apperror.NewAppError(err, "failed to decode request body in json", "", "DECODE_ERROR")
	}
	user, err := h.service.Create(context.Background(), project, h.storage)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) UpdateProject(w http.ResponseWriter, r *http.Request) error {
	var projectData Project
	err := json.NewDecoder(r.Body).Decode(&projectData)
	if err != nil {
		return apperror.NewAppError(err, "failed to decode request body in json", "", "DECODE_ERROR")
	}
	params := httprouter.ParamsFromContext(r.Context())
	projectData.ID = params.ByName("uuid")
	project, err := h.service.UpdateProject(context.Background(), h.storage, projectData)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
	return nil
}

func (h *handler) DeleteProject(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	oid := params.ByName("uuid")
	err := h.service.DeleteProject(context.Background(), h.storage, oid)
	if err != nil {
		return apperror.NewAppError(err, err.Error(), "", "")
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
