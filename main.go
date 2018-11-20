package main

import (
	"encoding/json"
	"healthagram-server/trans"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

var myServerURL = "http://localhost:8000"

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/test", testHTTP).Methods("GET")
	r.HandleFunc("/Html", htmlTestHTTP).Methods("GET")
	r.HandleFunc("/upload/{id}", imageUploadHTTP).Methods("POST")
	r.HandleFunc("/multiupload/{id}", imageMultiUploadHTTP).Methods("POST")
	r.HandleFunc("/json", trans.ContentCreateHTTP).Methods("POST")
	r.HandleFunc("/json/{id}/{page}", trans.ContentReadHTTP).Methods("GET")
	r.HandleFunc("/bulletin/image/{id}/{filename}", imageGetterHTTP).Methods("GET")
	http.ListenAndServe(":8000", r)
}
func imageGetterHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	image, err := os.Open("./image/" + id + "/" + filename)
	if err != nil {
		log.Print(err)
	}
	data := make([]byte, 32<<20)
	count, err := image.Read(data)
	if err != nil {
		log.Print(err)
	}
	log.Print(count)
	w.Write(data)

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

func imageUploadHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// sha256EncreptedID := sha256.New()
	// sha256EncreptedID.Write([]byte(id))
	// directory := hex.EncodeToString(sha256EncreptedID.Sum(nil))
	path := "./image/" + id

	file, header, err := r.FormFile("uplTheFile")
	if err != nil {
		log.Print(err)
	}

	filepath := path + "/" + header.Filename
	resultFilepaths := myServerURL + "/bulletin/image/" + id + "/" + header.Filename

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 777)
	}
	dst, err := os.Create(filepath)
	defer dst.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 777)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(resultFilepaths))
	if err != nil {
		log.Print(err)
	}
}

func imageMultiUploadHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// sha256EncreptedID := sha256.New()
	// sha256EncreptedID.Write([]byte(id))
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Print(err)
	}
	// directory := hex.EncodeToString(sha256EncreptedID.Sum(nil))
	path := "./image/" + id

	m := r.MultipartForm
	files := m.File["uplTheFile"]
	var resultFilepaths []string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 777)
	}
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filepath := path + "/" + files[i].Filename
		resultFilepaths[i] = filepath
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 777)
		}
		dst, err := os.Create(filepath)
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resultFilepaths)
	if err != nil {
		log.Print(err)
	}
}
