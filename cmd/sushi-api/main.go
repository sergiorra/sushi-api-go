package main

import (
	sushi "github.com/sergiorra/sushi-api-go/pkg"
	"github.com/sergiorra/sushi-api-go/pkg/adding"
	"github.com/sergiorra/sushi-api-go/pkg/getting"
	"github.com/sergiorra/sushi-api-go/pkg/modifying"
	"github.com/sergiorra/sushi-api-go/pkg/removing"
	"github.com/sergiorra/sushi-api-go/pkg/server"
	"github.com/sergiorra/sushi-api-go/pkg/storage/inmem"
	"log"
	"net/http"
)

func main() {
	var sushi map[string]sushi.Sushi
	repo := inmem.NewRepository(sushi)
	gS := getting.NewService(repo)
	aS := adding.NewService(repo)
	mS := modifying.NewService(repo)
	rS := removing.NewService(repo)

	s := server.New(gS, aS, mS, rS)
	log.Fatal(http.ListenAndServe(":8080", s.Router()))

}
