// cart.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/database"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/middleware"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/models"
)

// AddToCart adds an item to a user's cart.
func AddToCart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	requestedUserID := params["user_id"]
	userIDValue := r.Context().Value(middleware.ContextKeyUserID)

	currentUserID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusUnauthorized)
		return
	}

	if strconv.Itoa(currentUserID) != requestedUserID {
		http.Error(w, "Forbidden: You can only modify your own cart", http.StatusForbidden)
		return
	}

	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	item.UserID = currentUserID

	_, err := database.DB.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)", item.UserID, item.ProductID, item.Quantity)
	if err != nil {
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
}

func GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	requestedUserID := params["user_id"]

	userIDValue := r.Context().Value(middleware.ContextKeyUserID)
	currentUserID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusUnauthorized)
		return
	}

	if strconv.Itoa(currentUserID) != requestedUserID {
		http.Error(w, "Forbidden: You can only view your own cart", http.StatusForbidden)
		return
	}

	query := `SELECT c.id, c.user_id, c.product_id, c.quantity, p.name, p.price 
	          FROM cart_items c 
	          JOIN products p ON c.product_id = p.id 
	          WHERE c.user_id = $1`

	rows, err := database.DB.Query(query, currentUserID)
	if err != nil {
		http.Error(w, "Error fetching cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CartResponse struct {
		ID        int     `json:"id"`
		UserID    int     `json:"user_id"`
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Name      string  `json:"product_name"`
		Price     float64 `json:"product_price"`
	}

	var cart []CartResponse
	for rows.Next() {
		var item CartResponse
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.Name, &item.Price)
		if err != nil {
			http.Error(w, "Error scanning cart item", http.StatusInternalServerError)
			return
		}
		cart = append(cart, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing cart item ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
