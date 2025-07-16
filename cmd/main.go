package main

import (
	"MarketVK/internal/handler"
	"MarketVK/internal/repository"
	"log"
	"net/http"
)

func main() {
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	mux := handler.NewRouter()
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
