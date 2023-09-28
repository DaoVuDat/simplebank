package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"simple_bank/util"
	"testing"
)

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	testDb, err = pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
