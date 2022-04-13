package main

import (
	"backend/internal/user"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Create router")
	router := httprouter.New()
	log.Println("Register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	initial(router)
}

func initial(router *httprouter.Router) {
	log.Println("Start application")
	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server is listening 0.0.0.0:10000")
	log.Fatal(server.Serve(listener))
}
