package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Speaker struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	SocialLinks map[string]string `json:"social_links,omitempty"`
}

type Response struct {
	Title       string    `json:"title,omitempty"`
	Venue       string    `json:"venue,omitempty"`
	Time        string    `json:"time,omitempty"`
	Date        string    `json:"date,omitempty"`
	Description string    `json:"description,omitempty"`
	Speakers    []Speaker `json:"speakers,omitempty"`
	Links       []string  `json:"links,omitempty"`
}

var currentDevFestResponse Response

func init() {
	if file, err := os.Open("./devfest_current.json"); err != nil {
		log.Fatalln(err)
	} else {
		if bytes, err := io.ReadAll(file); err != nil {
			log.Fatalln(err)
		} else {
			if err := json.Unmarshal(bytes, &currentDevFestResponse); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(w, 200, currentDevFestResponse)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at port " + port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err.Error())
	}
}
