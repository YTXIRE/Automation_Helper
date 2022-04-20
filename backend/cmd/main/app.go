package main

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/project"
	projectDb "backend/internal/project/db"
	"backend/internal/user"
	userDb "backend/internal/user/db"
	"backend/pkg/client/mongodb"
	"backend/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	mongoDBClient, err := mongodb.NewClient(&mongodb.Config{
		Ctx:      context.Background(),
		Host:     cfg.MongoDB.Host,
		Port:     cfg.MongoDB.Port,
		Username: cfg.MongoDB.Username,
		Password: cfg.MongoDB.Password,
		Database: cfg.MongoDB.Database,
		AuthDB:   cfg.MongoDB.AuthDB,
	})

	if err != nil {
		panic(err)
	}

	logger.Info("Create user storage")
	userStorage := userDb.NewStorage(mongoDBClient, "users", logger)
	logger.Info("Create project storage")
	projectStorage := projectDb.NewStorage(mongoDBClient, "projects", logger)

	logger.Info("Register auth handler")
	authHandler := auth.NewHandler(logger, userStorage)
	authHandler.Register(router)

	logger.Info("Register user handler")
	userHandler := user.NewHandler(logger, userStorage)
	userHandler.Register(router)

	logger.Info("Register project handler")
	projectHandler := project.NewHandler(logger, projectStorage)
	projectHandler.Register(router)

	initial(router, logger, cfg)
}

func initial(router *httprouter.Router, logger *logging.Logger, cfg *config.Config) {
	logger.Info("Start application")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infof("Server is listening  %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
