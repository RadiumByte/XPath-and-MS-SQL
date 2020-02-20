package main

import (
	"sync"
	"time"

	"github.com/RadiumByte/XPath-and-MS-SQL/api"
	"github.com/RadiumByte/XPath-and-MS-SQL/app"
	"github.com/RadiumByte/XPath-and-MS-SQL/dal"
	"github.com/powerman/structlog"
)

var log = structlog.New()

func run(errc chan<- error) {
	time.Sleep(time.Second * 10)

	var mutex sync.Mutex

	// TODO: init DAL here for MS SQL
	db, err := dal.NewPostgresDAL("kaznachey", "RfpyfxtqAF", "localhost:5432")
	if err != nil {
		errc <- err
		return
	}

	application := app.NewApplication(db, errc)
	server := api.NewWebServer(application)

	server.Start(errc)
}

func main() {
	log.Info("Server is preparing to start...")

	errc := make(chan error)
	go run(errc)
	if err := <-errc; err != nil {
		log.Fatal(err)
	}
}
