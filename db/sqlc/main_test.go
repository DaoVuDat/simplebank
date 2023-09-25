package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "pgx/v5"
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error

	testDb, err = pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
