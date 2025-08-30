package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type serverStruct struct {
	clients  map[*websocket.Conn]bool
	logger   *slog.Logger
	serve    *http.Server
	upgrader websocket.Upgrader
}

var (
	mutex    sync.Mutex
	messChan = make(chan messStruct)
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := serverStruct{
		logger: logger,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		},
		clients: make(map[*websocket.Conn]bool),
	}

	go server.broadcast()

	mux := server.newRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Starting started", "port", port)
	err := http.ListenAndServe("0.0.0.0:"+port, mux)
	server.logger.Error(err.Error())
}
