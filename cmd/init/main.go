package main

import (
	"log"
	api "name/api/http"
	"name/internal/database"
	"os"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	app := api.InitializeNewServer()
	database.InitializeDatabase()

	if err := app.Listen(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
