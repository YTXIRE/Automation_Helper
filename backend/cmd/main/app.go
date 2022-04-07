package main

import (
	"backend/internal/user"
	"backend/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()
	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	initial(router, logger)
}

func initial(router *httprouter.Router, logger logging.Logger) {
	logger.Info("Start Application")
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Server is listening port 4000")
	logger.Fatal(server.Serve(listener))
}
