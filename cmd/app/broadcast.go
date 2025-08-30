package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func (s *serverStruct) broadcast() {
	for {
		messobj := <-messChan
		msg := fmt.Sprintf("%v : %v", messobj.user, string(messobj.message))
		for conn := range s.clients {
			if conn != messobj.conn {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					s.logger.Error(err.Error())
					mutex.Lock()
					delete(s.clients, conn)
					mutex.Unlock()
				}
			}
		}
	}
}
