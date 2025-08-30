package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (s serverStruct) connect(w http.ResponseWriter, r *http.Request) {
	user := r.Header["User"][0]
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Couldn't upgrade to websocket", http.StatusInternalServerError)
		s.logger.Error(err.Error())
	}
	mutex.Lock()
	s.clients[conn] = true
	mutex.Unlock()
	s.logger.Info("conn added")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.closeConn(conn, websocket.CloseAbnormalClosure, err.Error())
			return
		}
		if string(message) == "\\q" {
			s.closeConn(conn, websocket.CloseNormalClosure, "bye..")
			return
		}
		messObj := messStruct{message: message, conn: conn, user: user}
		messChan <- messObj
	}
}

func (s serverStruct) closeConn(conn *websocket.Conn, code int, reason string) {
	mutex.Lock()
	delete(s.clients, conn)
	mutex.Unlock()

	msg := websocket.FormatCloseMessage(code, reason)
	if err := conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
		s.logger.Error(err.Error())
	}

	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				s.logger.Info("Recieved close from peer")
			} else {
				s.logger.Warn(err.Error())
			}
			break
		}
	}
	s.logger.Info("Conn closed gracefully.")
}
