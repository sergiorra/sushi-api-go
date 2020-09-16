package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Sushi struct {
	ID          string `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Ingredients string `json:"ingredients"`
}

var sushis []Sushi

func getAllSushiHandler(w http.ResponseWriter, r *http.Request) {

}

func getSushiHandler(w http.ResponseWriter, r *http.Request) {

}

func createSushiHandler(w http.ResponseWriter, r *http.Request) {

}

func updateSushiHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteSushiHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// initialize router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/sushi", getAllSushiHandler).Methods("GET")
	router.HandleFunc("/sushi/{id}", getSushiHandler).Methods("GET")
	router.HandleFunc("/sushi", createSushiHandler).Methods("POST")
	router.HandleFunc("/sushi/{id}", updateSushiHandler).Methods("UPDATE")
	router.HandleFunc("/sushi/{id}", deleteSushiHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
