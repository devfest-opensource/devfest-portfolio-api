package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Person struct {
	Name        string            `json:"name,omitempty"`
	Designation string            `json:"designation,omitempty"`
	Track       string            `json:"track,omitempty"`
	SocialLinks map[string]string `json:"social_links,omitempty"`
	Image       string            `json:"image,omitempty"`
	Title       string            `json:"title,omitempty"`
}

type Response struct {
	Title            string   `json:"title,omitempty"`
	RegistrationLink string   `json:"registration_link,omitempty"`
	Description      string   `json:"description,omitempty"`
	Venue            []string `json:"venue,omitempty"`
	Time             string   `json:"time,omitempty"`
	Date             string   `json:"date,omitempty"`
	Image            string   `json:"image,omitempty"`
	Speakers         []Person `json:"speakers,omitempty"`
	Organizers       []Person `json:"organizers,omitempty"`
	Links            []string `json:"links,omitempty"`
}

var currentDevFestResponse Response
var previousDevFestResponses []Response

func init() {
	responseReferencePerFilename := map[string]any{
		"devfest_current.json": &currentDevFestResponse,
		"devfest_past.json":    &previousDevFestResponses,
	}

	for filename, responseReference := range responseReferencePerFilename {
		file, err := os.Open(filename)
		handleError(err)

		bytes, err := io.ReadAll(file)
		handleError(err)

		handleError(json.Unmarshal(bytes, responseReference))
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(w, 200, currentDevFestResponse)
	})

	router.HandleFunc("/past", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponse(w, 200, previousDevFestResponses)
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
	handleError(json.NewEncoder(w).Encode(data))
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
