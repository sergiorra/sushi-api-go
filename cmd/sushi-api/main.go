package main

import (
	"flag"
	"fmt"
	"github.com/sergiorra/sushi-api-go/pkg/log/logrus"
	"log"
	"net/http"
	"os"
	"strconv"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
	"github.com/sergiorra/sushi-api-go/pkg/adding"
	"github.com/sergiorra/sushi-api-go/pkg/getting"
	"github.com/sergiorra/sushi-api-go/pkg/modifying"
	"github.com/sergiorra/sushi-api-go/pkg/removing"
	"github.com/sergiorra/sushi-api-go/pkg/server"
	"github.com/sergiorra/sushi-api-go/pkg/storage/inmem"
)

func main() {

	var (
		hostName, _     = os.Hostname()
		defaultServerID = fmt.Sprintf("%s-%s", os.Getenv("SUSHIAPI_NAME"), hostName)
		defaultHost     = os.Getenv("SUSHIAPI_SERVER_HOST")
		defaultPort, _  = strconv.Atoi(os.Getenv("SUSHIAPI_SERVER_PORT"))
	)

	host := flag.String("host", defaultHost, "define host of the server")
	port := flag.Int("port", defaultPort, "define port of the server")
	serverID := flag.String("server-id", defaultServerID, "define server identifier")

	var sushis map[string]sushi.Sushi
	logger := logrus.NewLogger()

	repo := inmem.NewRepository(sushis)
	gS := getting.NewService(repo, logger)
	aS := adding.NewService(repo)
	mS := modifying.NewService(repo)
	rS := removing.NewService(repo)

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)

	s := server.New(*serverID, gS, aS, mS, rS)

	fmt.Println("The sushi server is on tap now:", httpAddr)
	log.Fatal(http.ListenAndServe(":8080", s.Router()))

}
