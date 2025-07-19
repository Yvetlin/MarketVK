package main

import (
	"MarketVK/internal/db"
	"MarketVK/internal/handler"
	"log"
	"net/http"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	mux := handler.NewRouter()
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
