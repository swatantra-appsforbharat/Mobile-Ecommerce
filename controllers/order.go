package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/database"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/models"
)

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	requestedUserID := params["user_id"]
	currentUserID := r.Context().Value("user_id").(int)

	if strconv.Itoa(currentUserID) != requestedUserID {
		http.Error(w, "Forbidden: You can only place orders for your own account", http.StatusForbidden)
		return
	}

	// Fetch cart items
	rows, err := database.DB.Query(`SELECT product_id, quantity FROM cart_items WHERE user_id = $1`, currentUserID)
	if err != nil {
		http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		item.UserID = currentUserID
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			http.Error(w, "Error reading cart data", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Create order
	var orderID int
	err = database.DB.QueryRow(
		`INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id`,
		currentUserID, "placed",
	).Scan(&orderID)
	if err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	// Insert order items
	for _, item := range items {
		_, err := database.DB.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			orderID, item.ProductID, item.Quantity,
		)
		if err != nil {
			http.Error(w, "Failed to insert order items", http.StatusInternalServerError)
			return
		}
	}

	// Clear cart
	_, err = database.DB.Exec(`DELETE FROM cart_items WHERE user_id = $1`, currentUserID)
	if err != nil {
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Order placed",
		"order_id": orderID,
	})
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderIDStr := params["order_id"]

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`UPDATE orders SET status = 'cancelled' WHERE id = $1`, orderID)
	if err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Order cancelled",
		"order_id": orderID,
	})
}
