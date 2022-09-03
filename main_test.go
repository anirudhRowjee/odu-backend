package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var a App

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS products (
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)`

// TODO refactor these two functions, put them in `db/utils`?
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

// testing!
func TestMain(m *testing.M) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could Not Load DotEnv [%s]", err)
	}
	log.Printf("username %s password %s db %s",
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	// here `a` is a global variable for the application we want to test
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	// basic DB Cleanup and Initialization
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}
