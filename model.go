package main

import (
	"database/sql"
)

// this file contains all the model information
type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// CREATE
func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

// READ
func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id = $1", p.ID).Scan(&p.Name, &p.Price)
}

// UPDATE
func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)
	return err
}

// DELETE
func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products where id=$1", p.ID)
	return err
}

// Standalone function, pagination-supported product list
func getProducts(db *sql.DB, start, count int) ([]product, error) {

	// fetch all the rows that exist
	// use Query for select statements, Exec for delete / update / etc
	rows, err := db.Query("SELECT id, name, price FROM products LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	// keep the rows alive and safe from GC until we have consumed them all
	defer rows.Close()

	products := []product{}
	for rows.Next() {
		// prepare an empty product
		var p product
		// dump the latest row into this product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		// add this product to all the products
		products = append(products, p)
	}

	return products, nil

}
