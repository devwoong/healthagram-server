package trans

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type contents struct {
	Seq         int                   `json:"seq"`
	UserID      int                   `json:"user_id"`
	Date        string                `json:"date"`
	Author      string                `json:"author"`
	Title       string                `json:"title"`
	IsAnonym    string                `json:"isAnonym"`
	NoticeScope string                `json:"notice_scope"`
	Editors     [](map[string]string) `json:"editors"`
	Images      [](map[string]string) `json:"images"`
	Videos      [](map[string]string) `json:"videos"`
	Kewords     []string              `json:"kewords"`
}
type contentArray struct {
	Contents []contents `json:"contents"`
}

// ContentCreateHTTP insert contents at DB
func ContentCreateHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	jsonContent := contents{}
	err := decoder.Decode(&jsonContent)
	if err != nil {
		log.Print()
		fmt.Print(err)
	}
	for _, editors := range jsonContent.Editors {
		for k, v := range editors {
			log.Print(k + " : " + v)
		}
	}

	for _, images := range jsonContent.Images {
		for k, v := range images {
			log.Print(k + " : " + v)
		}
	}

	for _, videos := range jsonContent.Videos {
		for k, v := range videos {
			log.Print(k + " : " + v)
		}
	}
	defer r.Body.Close()
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("healthagram").C("bulletins")

	var result = contents{}
	err = c.Find(nil).Select(bson.M{"seq": 1}).Sort("{seq:-1}").Limit(1).One(&result)
	if err != nil {
		jsonContent.Seq = 1
	} else {
		jsonContent.Seq = result.Seq + 1
	}
	err = c.Insert(jsonContent)
	if err != nil {
		log.Fatal(err)
	}
}

// ContentReadHTTP Read requeted content
func ContentReadHTTP(w http.ResponseWriter, r *http.Request) {
	const maxBulletin int = 10
	vars := mux.Vars(r)
	id := vars["id"]
	page, err := strconv.ParseInt(vars["page"], 10, 32)
	log.Print(id)
	log.Print(page)
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collections := session.DB("healthagram").C("bulletins")
	var result = []contents{}
	err = collections.Find(nil).Sort("{date: 1}").Skip(int(page) * maxBulletin).Limit(maxBulletin).All(&result)
	if err != nil {
		log.Print(err)
	}
	var JObjectResult = contentArray{Contents: result}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(JObjectResult)
	if err != nil {
		log.Print(err)
	}
}
