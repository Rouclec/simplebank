package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" //Must ADD!! for code to be able to communicate with database
	"github.com/rouclec/simplebank/api"
	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/util"
)


func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error parsing database config: ", err)
	}

	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	store := db.NewStore(pool)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
