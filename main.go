package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/harryng22/simplebank/api"
	db "github.com/harryng22/simplebank/db/sqlc"
	"github.com/harryng22/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start the server", err)
	}
}
