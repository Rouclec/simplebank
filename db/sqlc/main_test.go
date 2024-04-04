package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rouclec/simplebank/util"
)

var testQueries *Queries
var pool *pgxpool.Pool

func TestMain(m *testing.M) {

	var err error


	config := util.Config{
		DBSource:            "postgresql://postgres:changemeinprod%21@localhost:5432/simple_bank?sslmode=disable",
		MigrationURL:        "file://db/migration",
		ServerAddress:       "0.0.0:8080",
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute * 15,
		Domain:              "localhost",
	}

	pool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	testQueries = New(pool)

	defer pool.Close() // Close the pool after tests

	os.Exit(m.Run())
}
