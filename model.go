package main

import (
	"database/sql"
	"errors"
)

// this file contains all the model information
type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// CREATE
func (p *product) createProduct(db *sql.DB) error {
	return errors.New("Not Implemented!")
}

// READ
func (p *product) getProduct(db *sql.DB) error {
	return errors.New("Not Implemented!")
}

// UPDATE
func (p *product) updateProduct(db *sql.DB) error {
	return errors.New("Not Implemented!")
}

// DELETE
func (p *product) deleteProduct(db *sql.DB) error {
	return errors.New("Not Implemented!")
}

// Standalone function, pagination-supported product list
func getProducts(db *sql.DB, start, count int) ([]product, error) {
	return nil, errors.New("not implemented!")
}
