package main

import (
	"dynamodb-url-shortener/db"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("root url accessed", time.Now().Local().UTC())
	}).Methods(http.MethodGet)

	router.HandleFunc("/app/", createShortCode).Methods(http.MethodPost)
	router.HandleFunc("/app/{shortcode}", accessURLWithShortCode).Methods(http.MethodGet)

	log.Println("starting server.....")
	http.ListenAndServe(":8080", router)
}

func createShortCode(rw http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	url := string(b)
	log.Println("URL", url)

	shortCode, err := db.SaveURL(url)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(CreateShortCodeResponse{ShortCode: shortCode})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	log.Println("short code for", url, shortCode)
}

type CreateShortCodeResponse struct {
	ShortCode string
}

func accessURLWithShortCode(rw http.ResponseWriter, req *http.Request) {
	shortCode := mux.Vars(req)["shortcode"]

	url, err := db.GetLongURL(shortCode)

	if err != nil {
		if errors.Is(err, db.ErrUrlNotFound) {
			http.Error(rw, err.Error(), http.StatusNotFound)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusFound)

	log.Println("found short code for", url, shortCode)
}
