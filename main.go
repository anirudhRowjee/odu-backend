package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could Not Load DotEnv [%s]", err)
	}

	// declare the App
	a := App{}

	// Start off the app with some Environment Variables
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	// Run on port 8010
	a.Run(":8010")
}
