package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Sushi struct {
	ID          string `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Ingredients string `json:"ingredients"`
}

var sushis []Sushi

func main() {
	// initialize router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/sushi", getSushis).Methods("GET")
	router.HandleFunc("/sushi/{id}", getSushi).Methods("GET")
	router.HandleFunc("/sushi", createSushi).Methods("POST")
	router.HandleFunc("/sushi/{id}", updateSushi).Methods("UPDATE")
	router.HandleFunc("/sushi/{id}", deleteSushi).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
