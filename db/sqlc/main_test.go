package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rouclec/simplebank/util"
)

var testQueries *Queries
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error parsing database config: ", err)
	}

	pool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	testQueries = New(pool)

	defer pool.Close() // Close the pool after tests

	os.Exit(m.Run())
}
