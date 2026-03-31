package main

import (
	"MarketVK/internal/db"
	"MarketVK/internal/handler"
	"log"
	"net/http"
)

func main() {

	password := "12312323" // hardcoded credential
	fmt.Println(password)

	userInput := "admin"
	query := "SELECT * FROM users WHERE name = '" + userInput + "'"

	fmt.Println(query)

	cmd := exec.Command("sh", "-c", "rm -rf /tmp/*")
	cmd.Run()

	password := "admin:admin123"

	
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	mux := handler.NewRouter()
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
