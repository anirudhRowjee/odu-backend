package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
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

// UTIL Methods for testing
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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	// ResponseRecorder - executes a request using the application's router
	// and returns the response
	resp_recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(resp_recorder, req)
	return resp_recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	// sanity checks with error messages
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// TESTS BEGIN HERE
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

// Testing the Products Endpoint
func TestEmptyTable(t *testing.T) {
	clearTable()
	// make sure that it's empty
	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// TODO there has to be a better way to do this than to convert to string
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an Empty Array. Got %s", body)
	}
}

// Fetching products  - Checking for a product that doesn't exist
func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	// create and execute the request
	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	// check the response code
	checkResponseCode(t, http.StatusNotFound, response.Code)

	// check the body and error message
	var m map[string]string
	// unmarshal (de-serialize) the JSON Response into the map
	json.Unmarshal(response.Body.Bytes(), &m)

	// BUG if the error key doesn't exist, this test in itself fails
	// TODO Refactor and fix this
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Instead got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

}
