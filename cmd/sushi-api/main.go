package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	sushi "github.com/sergiorra/sushi-api-go/pkg"
	"github.com/sergiorra/sushi-api-go/pkg/adding"
	"github.com/sergiorra/sushi-api-go/pkg/getting"
	"github.com/sergiorra/sushi-api-go/pkg/log/logrus"
	"github.com/sergiorra/sushi-api-go/pkg/modifying"
	"github.com/sergiorra/sushi-api-go/pkg/removing"
	"github.com/sergiorra/sushi-api-go/pkg/server"
	"github.com/sergiorra/sushi-api-go/pkg/storage/cockroach"
	"github.com/sergiorra/sushi-api-go/pkg/storage/inmem"
	"github.com/sergiorra/sushi-api-go/pkg/storage/mysql"
)

func main() {

	var (
		hostName, _     = os.Hostname()
		defaultServerID = fmt.Sprintf("%s-%s", os.Getenv("SUSHIAPI_NAME"), hostName)
		defaultHost     = os.Getenv("SUSHIAPI_SERVER_HOST")
		defaultPort, _  = strconv.Atoi(os.Getenv("SUSHIAPI_SERVER_PORT"))
		defaultDB		= "inmem"
	)

	host := flag.String("host", defaultHost, "define host of the server")
	port := flag.Int("port", defaultPort, "define port of the server")
	serverID := flag.String("server-id", defaultServerID, "define server identifier")
	database := flag.String("database", defaultDB, "initialize the api using the given db engine")
	flag.Parse()

	var sushis map[string]sushi.Sushi
	logger := logrus.NewLogger()

	repo := initializeRepo(database, sushis)
	gS := getting.NewService(repo, logger)
	aS := adding.NewService(repo)
	mS := modifying.NewService(repo)
	rS := removing.NewService(repo)

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)

	s := server.New(*serverID, gS, aS, mS, rS)

	fmt.Println("The sushi server is on tap now:", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, s.Router()))

}

func initializeRepo(database *string, sushis map[string]sushi.Sushi) sushi.Repository {
	var repo sushi.Repository
	switch *database {
	case "cockroach":
		repo = newCockroachRepository()
	case "mysql":
		repo = newMySQLRepository()
	default:
		repo = inmem.NewRepository(sushis)
	}
	return repo
}

func newCockroachRepository() sushi.Repository {
	cockroachAddr := os.Getenv("COCKROACH_ADDR")
	cockroachDBName := os.Getenv("COCKROACH_DB")

	cockroachConn, err := cockroach.NewConn(cockroachAddr, cockroachDBName)
	if err != nil {
		log.Fatal(err)
	}
	return cockroach.NewRepository(cockroachConn)
}

func newMySQLRepository() sushi.Repository {
	mysqlAddr := os.Getenv("MYSQL_ADDR")
	mysqlDBName := os.Getenv("MYSQL_DB")

	mysqlConn, err := mysql.NewConn(mysqlAddr, mysqlDBName)
	if err != nil {
		log.Fatal(err)
	}
	return mysql.NewRepository("gophers", mysqlConn)
}