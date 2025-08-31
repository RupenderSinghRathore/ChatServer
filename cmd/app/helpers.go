package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type messStruct struct {
	user    string
	message []byte
	conn    *websocket.Conn
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Panic occured", http.StatusInternalServerError)
				app.logger.Error(fmt.Sprintf("Panic occured: %s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
