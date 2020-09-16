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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sushis)
}

func getSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range sushis {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newSushi := Sushi{}
	json.NewDecoder(r.Body).Decode(&newSushi)
	newSushi.ID = strconv.Itoa(len(sushis)+1)
	sushis = append(sushis, newSushi)
	json.NewEncoder(w).Encode(newSushi)
}

func updateSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range sushis {
		if item.ID == params["id"] {
			sushis = append(sushis[:i], sushis[i+1:]...)
			newSushi := Sushi{}
			json.NewDecoder(r.Body).Decode(&newSushi)
			newSushi.ID = params["id"]
			sushis = append(sushis, newSushi)
			json.NewEncoder(w).Encode(newSushi)
			return
		}
	}
}

func deleteSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range sushis {
		if item.ID == params["id"] {
			sushis = append(sushis[:i], sushis[i+1:]...)
			return
		}
	}
	json.NewEncoder(w).Encode(sushis)
}

func main() {
	// initialize router
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/sushi", getAllSushiHandler).Methods("GET")
	router.HandleFunc("/sushi/{id}", getSushiHandler).Methods("GET")
	router.HandleFunc("/sushi", createSushiHandler).Methods("POST")
	router.HandleFunc("/sushi/{id}", updateSushiHandler).Methods("PUT")
	router.HandleFunc("/sushi/{id}", deleteSushiHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
