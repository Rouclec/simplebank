package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" //Must ADD!! for code to be able to communicate with database
	"github.com/rouclec/simplebank/api"
	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/gapi"
	"github.com/rouclec/simplebank/pb"
	"github.com/rouclec/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(config, store)
	// runGinServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("error creating server: ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("GRPC Address: %s", config.GRPCServerAddress)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("Error creating gRPC listener: ", err)
	}

	log.Printf("gRPC listener running on port %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Error starting gRPC server: ", err)
	}

	log.Printf("gRPC server running on %s", config.GRPCServerAddress)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("error creating server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
