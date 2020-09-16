package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Sushi struct {
	ID          int `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Ingredients []string `json:"ingredients"`
}

var sushiData []Sushi

func getAllSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sushiData)
}

func getSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, item := range sushiData {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newSushi := Sushi{}
	json.NewDecoder(r.Body).Decode(&newSushi)
	newSushi.ID = len(sushiData)+1
	sushiData = append(sushiData, newSushi)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSushi)
}

func updateSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i, item := range sushiData {
		if item.ID == id {
			newSushi := Sushi{}
			json.NewDecoder(r.Body).Decode(&newSushi)
			newSushi.ID = id
			sushiData[i] = newSushi
			json.NewEncoder(w).Encode(newSushi)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteSushiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i, item := range sushiData {
		if item.ID == id{
			sushiData = append(sushiData[:i], sushiData[i+1:]...)
			json.NewEncoder(w).Encode(sushiData)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
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
