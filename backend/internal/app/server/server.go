package server

import (
	"backend/internal/app/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/users", s.GetUsers())
}

func (s *Server) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(s.logger); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *Server) GetUsers() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email, err := s.store.User().FindByEmail("admin@it-paradise.com")
		if err != nil {
			logrus.Error(err)
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		marshal, err := json.Marshal(email)
		if err != nil {
			return
		}
		writer.Write(marshal)
	}
}
