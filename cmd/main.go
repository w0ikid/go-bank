package main

import (
	"log"
	"context"
	"github.com/w0ikid/go-bank/api"
	db "github.com/w0ikid/go-bank/db/sqlc"
	"github.com/joho/godotenv"
	"github.com/w0ikid/go-bank/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env не найден или не загружен")
	}

	// Load configuration
	loader := pkg.CleanenvLoader{}
	
	cfg := pkg.InitConfig(loader, "config.yaml")

	
	// Initialize database connection
	conn, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	// Initialize store and server

	store := db.NewStore(conn)
	server := api.NewServer(store)
	
	if err := server.Start(cfg.Server.Address); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}