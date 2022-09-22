package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

func addProducts(count int) {
	// add a standard number of products, dummy data
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec(
			"INSERT INTO products(name, price) VALUES ($1, $2)",
			"Product "+strconv.Itoa(i),
			(i+1.0)*10,
		)
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

	clearTable()

	// Create a new JSON Object to insert into the database, and then create the
	// Request to insert it into the DB

	// Declaring the JSON Object as a byte array (ASCII stonks?)
	var jsonStr = []byte(`{"name": "test product", "price": 11.22}`)

	// converting bytes.NewBuffer lets this request take ownership of the bytearray
	// converting it to an I/O Compliant Buffer
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Run this request
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	// dump returned JSON Object into a map[string]interface{}
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	// Sanity Check
	// TODO There has to be a better way of doing this
	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}
	if m["price"] != 11.22 {
		t.Errorf("Expected price to be '11.22'. Got '%v'", m["price"])
	}
	// unmarshalling converts to float when nothing else is specified
	if m["id"] != 1.00 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}

}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	// get the product
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {

	// Cleanup and Sanity Checks
	clearTable()
	addProducts(1)

	// Get the first product
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	// unmarshal the response into originalProduct
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)
	log.Println("Checked response 1")

	var jsonStr = []byte(`{"name": "Updated Product Name", "price": 10000}`)
	log.Println("Modifications made")

	// create a new request with the unmarshalled object
	update_req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	update_req.Header.Set("Content-Type", "application/json")
	log.Println("Created a new request")

	// run update request
	response_updated := executeRequest(update_req)

	//  ensure update is successful, check return body of response_updated
	var newProduct map[string]interface{}
	json.Unmarshal(response_updated.Body.Bytes(), &newProduct)

	// check to see if our values have been updated
	// TODO see if we can use a JSON Diff library of some sort
	if newProduct["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], newProduct["id"])
	}

	if newProduct["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], newProduct["name"], newProduct["name"])
	}

	if newProduct["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], newProduct["price"], newProduct["price"])
	}

}

func TestDeleteProduct(t *testing.T) {

	// setup
	clearTable()
	addProducts(1)

	// check if the product got created
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// delete the product
	req_delete, _ := http.NewRequest("DELETE", "/product/1", nil)
	response_delete := executeRequest(req_delete)
	checkResponseCode(t, http.StatusOK, response_delete.Code)

	// check if the product fetch again fails
	req_check, _ := http.NewRequest("GET", "/product/1", nil)
	response_check := executeRequest(req_check)
	checkResponseCode(t, http.StatusNotFound, response_check.Code)
}
