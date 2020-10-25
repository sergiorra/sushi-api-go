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

type server struct {
	serverID string
	httpAddr string

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

func New(serverID string, gS getting.Service, aS adding.Service, mS modifying.Service, rS removing.Service) Server {
	a := &server{serverID: serverID, getting: gS, adding: aS, modifying: mS, removing: rS}
	router(a)
	return a
}

func router(s *server) {
	r := mux.NewRouter()

	r.Use(newServerMiddleware(s.serverID))

	r.HandleFunc("/sushi", s.GetSushis).Methods(http.MethodGet)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", s.GetSushi).Methods(http.MethodGet)
	r.HandleFunc("/sushi", s.AddSushi).Methods(http.MethodPost)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", s.ModifySushi).Methods(http.MethodPut)
	r.HandleFunc("/sushi/{ID:[a-zA-Z0-9_]+}", s.RemoveSushi).Methods(http.MethodDelete)

	s.router = r
}

func (s *server) Router() http.Handler {
	return s.router
}

func (s *server) GetSushis(w http.ResponseWriter, r *http.Request) {
	sushis, _ := s.getting.GetSushis(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sushis)
}

func (s *server) GetSushi(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sushi := s.getting.GetSushiByID(r.Context(), params["ID"])
	w.Header().Set("Content-Type", "application/json")
	if sushi == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode("Sushi Not found")
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
func (s *server) AddSushi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var sushi addSushiRequest
	err := decoder.Decode(&sushi)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}

	if err := s.adding.AddSushi(r.Context(), sushi.ID, sushi.ImageNumber, sushi.Name, sushi.Ingredients); err != nil {
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
func (s *server) ModifySushi(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var sushi modifySushiRequest
	err := decoder.Decode(&sushi)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error unmarshalling request body")
		return
	}
	vars := mux.Vars(r)
	if err := s.modifying.ModifySushi(r.Context(), vars["ID"], sushi.ImageNumber, sushi.Name, sushi.Ingredients); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can't modify a sushi")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveSushi remove a sushi
func (s *server) RemoveSushi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s.removing.RemoveSushi(r.Context(), vars["ID"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
