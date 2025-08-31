package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func (app *application) broadcast() {
	for {
		messobj := <-messChan
		msg := fmt.Sprintf("%v : %v", messobj.user, string(messobj.message))
		for conn := range app.clients {
			if conn != messobj.conn {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					app.logger.Error(err.Error())
					mutex.Lock()
					delete(app.clients, conn)
					mutex.Unlock()
				}
			}
		}
	}
}
