package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", testHTTP).Methods("GET")
	r.HandleFunc("/Html", htmlTestHTTP).Methods("GET")
	r.HandleFunc("/upload", imageGetterHTTP).Methods("POST")
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

	// r.ParseMultipartForm(32 << 20)
	// file, handler, err := r.FormFile("uplTheFile")
	// if err != nil {
	// 	log.Print("Handler or file error : " + err.Error())
	// 	return
	// }
	// defer file.Close()
	// fmt.Fprintf(w, "%v", handler.Header)
	// f, err := os.OpenFile("./image/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Print("save file error : " + err.Error())
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, file)

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
