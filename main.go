package main

import (
	"database/sql"
	"log"

	"github.com/NatdanaiKhe/simplebank/api"
	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DB_Driver, config.DB_URL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.SERVER_ADDRESS)

	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
