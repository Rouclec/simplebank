package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *Queries
var pool *pgxpool.Pool

const dbSource = "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable"

func TestMain(m *testing.M) {
	var err error

	// config, err := util.LoadConfig("../..")
	// if err != nil {
	// 	log.Fatal("Error parsing database config: ", err)
	// }

	pool, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	testQueries = New(pool)

	defer pool.Close() // Close the pool after tests

	os.Exit(m.Run())
}
