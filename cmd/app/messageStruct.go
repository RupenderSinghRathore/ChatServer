package main

import "github.com/gorilla/websocket"

type messStruct struct {
	user    string
	message []byte
	conn    *websocket.Conn
}
