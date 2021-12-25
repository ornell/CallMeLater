package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	conf "github.com/tkanos/gonfig"
	"os"
)

var (
	storage  RequestStorage
	requests = make(chan *requestData)
	hasMore  = false
)

type Configuration struct {
	dbType           string `json:"db_type" env:"DBTYPE"`
	connectionString string `json:"connection_string" env:"CONSTRING"`
}

func main() {
	C := Configuration{}
	err := conf.GetConf("appsettings.json", &C)
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if C.dbType == "PSQL" {
		storage = newPgStorage(C.connectionString)
	} else {
		storage = newSqLiteStorage(C.connectionString)
	}

	go consumeLoop()
	handleRequests()
}
