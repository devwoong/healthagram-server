package main

import (
	"log"
	"net"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", testHTTP).Methods("GET")
	r.HandleFunc("/Html", htmlTestHTTP).Methods("GET")
	//	server := &http.Server{
	//		Handler: r,
	//		Addr:    "localhost:8000",
	//	}
	//	server.ListenAndServe()
	http.ListenAndServe(":8000", r)
}

func testHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))

	ip, port, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Print("err ")
	}
	log.Print("request from other user IP : " + ip + " port : " + port)
}

func htmlTestHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("test.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, nil)
}
