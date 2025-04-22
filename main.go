package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/olaniyi38/BE/api"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/util"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("error loading config files %v", err)
		return
	}

	DB_URL := config.DBSource

	conn, err := sql.Open(config.DBDriver, DB_URL)

	if err != nil {
		log.Fatalf("error opening DB connection: %v \n", err)
		return
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("error creating server %v \n", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatalf("unable to start serve: %v \n", err)
	}

}

//tables needed for a bank
//accounts - user data
//entries - records of money moving in and out of any account
//transfers - records of money moving between accounts
