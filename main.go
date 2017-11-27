package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type test_s struct {
	id    string
	os    string
	phone string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", testHTTP).Methods("GET")
	r.HandleFunc("/Html", htmlTestHTTP).Methods("GET")
	r.HandleFunc("/upload", imageGetterHTTP).Methods("POST")
	r.HandleFunc("/json", jsonGetterHTTP).Methods("POST")
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

func imageGetterHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Print(err)
	}
	m := r.MultipartForm
	files := m.File["uplTheFile"]
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dst, err := os.Create("./image/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
func jsonGetterHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonContent test_s
	var jsondata map[string]string
	err := decoder.Decode(&jsondata)
	if err != nil {
		log.Print(err)
	}

	for k, v := range jsondata {
		log.Print(k + " : " + v)
	}
	defer r.Body.Close()
	log.Print(jsonContent.id)

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("healthagram").C("bulletins")
	err = c.Insert(jsondata)
	if err != nil {
		log.Fatal(err)
	}
}
