package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sergiorra/sushi-api-go/pkg/adding"
	"github.com/sergiorra/sushi-api-go/pkg/getting"
	"github.com/sergiorra/sushi-api-go/pkg/modifying"
	"github.com/sergiorra/sushi-api-go/pkg/removing"
)

type api struct {
	router http.Handler
	getting  getting.Service
	modifying  modifying.Service
	adding  adding.Service
	removing  removing.Service
}

type Server interface {
	Router() http.Handler
	GetSushis(w http.ResponseWriter, r *http.Request)
	GetSushi(w http.ResponseWriter, r *http.Request)
	AddSushi(w http.ResponseWriter, r *http.Request)
	ModifySushi(w http.ResponseWriter, r *http.Request)
	RemoveSushi(w http.ResponseWriter, r *http.Request)
}

func New(gS getting.Service, aS adding.Service, mS modifying.Service, rS removing.Service) Server {
	a := &api{getting: gS, adding: aS, modifying: mS, removing: rS}
	r := mux.NewRouter()
	r.HandleFunc("/sushi", a.GetSushis).Methods(http.MethodGet)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", a.GetSushi).Methods(http.MethodGet)
	r.HandleFunc("/sushi", a.AddSushi).Methods(http.MethodPost)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", a.ModifySushi).Methods(http.MethodPut)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", a.RemoveSushi).Methods(http.MethodDelete)
	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

func (a *api) GetSushis(w http.ResponseWriter, r *http.Request) {
	sushis, _ := a.getting.GetSushis(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sushis)
}

func (a *api) GetSushi(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sushi, err := a.getting.GetSushiByID(r.Context(), params["ID"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // We use not found for simplicity
		json.NewEncoder(w).Encode("Sushi Not found")
		return
	}

	json.NewEncoder(w).Encode(sushi)
}

type addSushiRequest struct {
	ID    			string 		`json:"id"`
	ImageNumber  	string 		`json:"imageNumber"`
	Name 			string 		`json:"name"`
	Ingredients   	[]string  	`json:"ingredients"`
}

// AddSushi save a sushi
func (a *api) AddSushi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var s addSushiRequest
	err := decoder.Decode(&s)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}

	if err := a.adding.AddSushi(r.Context(), s.ID, s.ImageNumber, s.Name, s.Ingredients); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can't create a sushi")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type modifySushiRequest struct {
	ImageNumber  	string 		`json:"imageNumber"`
	Name 			string 		`json:"name"`
	Ingredients   	[]string  	`json:"ingredients"`
}

// ModifySushi modify sushi data
func (a *api) ModifySushi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var s modifySushiRequest
	err := decoder.Decode(&s)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}
	vars := mux.Vars(r)
	if err := a.modifying.ModifySushi(r.Context(), vars["ID"], s.ImageNumber, s.Name, s.Ingredients); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can't modify a sushi")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveSushi remove a sushi
func (a *api) RemoveSushi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a.removing.RemoveSushi(r.Context(), vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
