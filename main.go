package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)

	server, err := api.NewServer(*config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	server.Start(config.ServerAddress)
}
