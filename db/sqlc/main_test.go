package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://doniback:secret@localhost:5433/go_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool


func TestMain(m *testing.M) {
	var err error
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}