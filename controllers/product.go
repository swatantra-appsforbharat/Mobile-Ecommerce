// product.go
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/swatantra-appsforbharat/mobile-commercebackend/database"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/models"
)

// CreateProduct adds a new product to the database.
func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`
	err = database.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Quantity).Scan(&product.ID)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		fmt.Println("Insert Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Check Redis cache
	cacheKey := "all_products"
	cached, err := database.RDB.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		// Return cached response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	// Fetch from database
	rows, err := database.DB.Query(`SELECT id, name, description, price, quantity FROM products`)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	// Convert to JSON
	responseJSON, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
		return
	}

	// Store in Redis for 5 minutes
	database.RDB.Set(database.Ctx, cacheKey, responseJSON, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// GetProductByID fetches a product by its ID.
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	err := database.DB.QueryRow("SELECT id, name, description, price, quantity FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity)

	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
