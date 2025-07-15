package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/w0ikid/go-bank/api"
	db "github.com/w0ikid/go-bank/db/sqlc"
	"github.com/w0ikid/go-bank/util"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env не найден или не загружен")
	}

	cfg := util.InitConfig(util.CleanenvLoader{}, "config.yaml")

	// Initialize database connection
	conn, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	log.Println("Connected to database")

	// Initialize store and server

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(cfg.Server.Address); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
