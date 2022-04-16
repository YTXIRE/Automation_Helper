package project

import (
	"backend/internal/apperror"
	"backend/internal/handlers"
	"backend/pkg/logging"
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
	router.HandlerFunc(http.MethodGet, projectsUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, projectURL, apperror.Middleware(h.GetProjectByUUID))
	router.HandlerFunc(http.MethodPost, projectsUrl, apperror.Middleware(h.CreateProject))
	router.HandlerFunc(http.MethodPut, projectURL, apperror.Middleware(h.UpdateProject))
	router.HandlerFunc(http.MethodDelete, projectURL, apperror.Middleware(h.DeleteProject))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("GetList"))
	return nil
}

func (h *handler) GetProjectByUUID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) CreateProject(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) UpdateProject(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) DeleteProject(w http.ResponseWriter, r *http.Request) error {
	return nil
}
