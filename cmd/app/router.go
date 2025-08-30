package main

import "net/http"

func (s serverStruct) newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/connect", s.connect)

	return mux
}
