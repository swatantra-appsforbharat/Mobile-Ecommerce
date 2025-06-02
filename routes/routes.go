package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/controllers"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/middleware"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Product routes
	r.Handle("/products", middleware.JWTMiddleware(http.HandlerFunc(controllers.AddProduct))).Methods("POST")
	r.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", controllers.GetProductByID).Methods("GET")

	// Cart routes
	r.Handle("/cart/{user_id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.AddToCart))).Methods("POST")
	r.Handle("/cart/{user_id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetCartByUserID))).Methods("GET")
	r.Handle("/cart/item/{id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.RemoveFromCart))).Methods("DELETE")

	// Order routes
	r.Handle("/orders/{user_id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.PlaceOrder))).Methods("POST")
	r.Handle("/Cancelorder/{order_id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.CancelOrder))).Methods("DELETE")
	// User routes
	r.Handle("/users/{id}", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetUserByID))).Methods("GET")
	r.Handle("/users", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetAllUsers))).Methods("GET")

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"pong"}`))
	}).Methods("GET")

	return r
}
