package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/poohda-go/cmd/api"
	"github.com/poohda-go/db"
	"github.com/poohda-go/store"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error with the logger: %s", err.Error())
	}

	newDb, err := db.NewPostgresDb()
	if err != nil {
		log.Fatalf("Error with the logger: %s", err.Error())
	}

	zapLogger := logger.Sugar()
	store := store.NewStore(newDb)
	db.InitializeDb(newDb)
	db.AddMigrations(newDb)

	defer newDb.Close()

	server := api.NewApplication(zapLogger, store)

	if err := server.Run(); err != nil {
		log.Printf("Error Running: %s", err.Error())
	}
}
