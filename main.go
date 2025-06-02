package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/database"
	"github.com/swatantra-appsforbharat/mobile-commercebackend/routes"
)

func main() {
	database.Connect()
	database.InitRedis() // ← Add this

	router := routes.RegisterRoutes()
	handler := cors.Default().Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
