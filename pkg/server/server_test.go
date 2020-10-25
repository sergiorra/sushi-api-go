package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sergiorra/sushi-api-go/pkg/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergiorra/sushi-api-go/pkg/adding"
	"github.com/sergiorra/sushi-api-go/pkg/getting"
	"github.com/sergiorra/sushi-api-go/pkg/modifying"
	"github.com/sergiorra/sushi-api-go/pkg/removing"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
	"github.com/sergiorra/sushi-api-go/pkg/storage/inmem"

	"github.com/sergiorra/sushi-api-go/cmd/sample-data"
)

func TestGetSushis(t *testing.T) {
	req, err := http.NewRequest("GET", "/sushi", nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}

	s := buildServer()
	resRecorder := httptest.NewRecorder()

	s.GetSushis(resRecorder, req)

	res := resRecorder.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var got []sushi.Sushi
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("could not unmarshall response %v", err)
	}

	expected := len(sample.Sushis)

	if len(got) != expected {
		t.Errorf("expected %v sushis, got: %v sushi", sample.Sushis, got)
	}
}

func TestGetSushi(t *testing.T) {

	testData := []struct {
		name   string
		s      *sushi.Sushi
		status int
		err    string
	}{
		{name: "sushi found", s: sushiSample(), status: http.StatusOK},
		{name: "sushi not found", s: &sushi.Sushi{ID: "123"}, status: http.StatusNotFound, err: "Sushi Not found"},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/sushi/%s", tt.s.ID)
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			s := buildServer()

			resRecorder := httptest.NewRecorder()
			s.Router().ServeHTTP(resRecorder, req)

			res := resRecorder.Result()

			defer res.Body.Close()
			if tt.status != res.StatusCode {
				t.Errorf("expected %d, got: %d", tt.status, res.StatusCode)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tt.err == "" {
				var got *sushi.Sushi
				err = json.Unmarshal(b, &got)
				if err != nil {
					t.Fatalf("could not unmarshall response %v", err)
				}

				if got.ID != tt.s.ID {
					t.Fatalf("expected %v, got: %v", *tt.s, *got)
				}
			}
		})
	}

}

func TestAddSushi(t *testing.T) {
	bodyJSON := []byte(`{
        "ID": "01D3XZ38GYT",
        "imageNumber": "4",
        "name": "Dragon Roll",
        "ingredients": ["Crab", "Cucumber", "Avocado", "Eel"]
    }`)
	req, err := http.NewRequest("POST", "/sushi", bytes.NewBuffer(bodyJSON))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	s := buildServer()
	resRecorder := httptest.NewRecorder()

	s.AddSushi(resRecorder, req)
	res := resRecorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected %d, got: %d", http.StatusCreated, res.StatusCode)
	}
}

func TestModifySushi(t *testing.T) {
	bodyJSON := []byte(`{
        "imageNumber": "4",
        "name": "Dragon Roll",
		"ingredients": ["Crab", "Cucumber", "Avocado", "Eel"]
    }`)
	req, err := http.NewRequest("PUT", "/sushi/01D3XZ38KLE", bytes.NewBuffer(bodyJSON))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	s := buildServer()
	resRecorder := httptest.NewRecorder()

	s.ModifySushi(resRecorder, req)
	res := resRecorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected %d, got: %d", http.StatusNoContent, res.StatusCode)
	}
}

func TestRemoveSushi(t *testing.T) {

	uri := fmt.Sprintf("/sushi/%s", "01D3XZ38KLE")
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}

	s := buildServer()

	resRecorder := httptest.NewRecorder()
	s.Router().ServeHTTP(resRecorder, req)

	res := resRecorder.Result()

	defer res.Body.Close()
	if http.StatusNoContent != res.StatusCode {
		t.Errorf("expected %d, got: %d", http.StatusNoContent, res.StatusCode)
	}

}

func sushiSample() *sushi.Sushi {
	return &sushi.Sushi{
		ID:    "01D3XZ38KDR",
		ImageNumber:  "1",
		Name:   "California Roll",
		Ingredients: []string {"Crab", "Avocado", "Cucumber", "Sesame seeds"},
	}
}

func buildServer() Server {
	repo := inmem.NewRepository(sample.Sushis)
	fetching := getting.NewService(repo, log.NewNoopLogger())
	adding := adding.NewService(repo)
	modifying := modifying.NewService(repo)
	removing := removing.NewService(repo)

	return New("test", fetching, adding, modifying, removing)
}
