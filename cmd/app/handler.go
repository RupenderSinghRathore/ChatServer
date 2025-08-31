package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

func (app application) connect(w http.ResponseWriter, r *http.Request) {
	users, ok := r.Header["User"]
	if !ok || len(users) < 1 {
		app.serverError(w, "Bad Request")
		return
	}
	user := users[0]
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.serverError(w, err.Error())
	}
	mutex.Lock()
	app.clients[conn] = true
	mutex.Unlock()
	app.logger.Info("conn added")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			app.closeConn(conn, websocket.CloseAbnormalClosure, err.Error())
			return
		}
		if string(message) == "\\q" {
			app.closeConn(conn, websocket.CloseNormalClosure, "bye..")
			return
		}
		messObj := messStruct{message: message, conn: conn, user: user}
		messChan <- messObj
	}
}

func (app application) closeConn(conn *websocket.Conn, code int, reason string) {
	mutex.Lock()
	delete(app.clients, conn)
	mutex.Unlock()

	msg := websocket.FormatCloseMessage(code, reason)
	if err := conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
		app.logger.Error(err.Error())
	}

	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				app.logger.Info("Recieved close from peer")
			} else {
				app.logger.Warn(err.Error())
			}
			break
		}
	}
	app.logger.Info("Conn closed gracefully.")
}

func (app application) serverError(w http.ResponseWriter, errString string) {
	http.Error(w, "NOT FOUND", http.StatusInternalServerError)
	app.logger.Error(errString)
	fmt.Fprintf(os.Stdout, "trace: %s\n", string(debug.Stack()))
}
