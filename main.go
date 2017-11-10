package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", testHTTP).Methods("GET")
	server := &http.Server{
		Handler: r,
		Addr:    "localhost:8000",
	}
	server.ListenAndServe()
}

func testHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	log.Print("request from other user")
}
