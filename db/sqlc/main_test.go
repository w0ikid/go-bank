package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/w0ikid/go-bank/util"
	"log"
)

var testStore Store

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(".env не найден или не загружен")
	}

	cfg := util.InitConfig(util.CleanenvLoader{}, "../../config.yaml")

	conn, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
