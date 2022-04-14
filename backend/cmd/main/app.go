package main

import (
	"backend/internal/config"
	"backend/internal/user"
	"backend/pkg/logging"
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

	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
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
