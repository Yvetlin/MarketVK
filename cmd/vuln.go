package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

// 🔥 HARDCODED SECRET (gosec G101)
var password = "super-secret-password"

// 🔥 SQL INJECTION
func buildQuery(userInput string) string {
	return "SELECT * FROM users WHERE name = '" + userInput + "'"
}

// 🔥 COMMAND INJECTION
func runCommand(input string) {
	cmd := exec.Command("sh", "-c", "echo "+input)
	cmd.Run() // G104 - error not handled
}

// 🔥 HTTP WITHOUT TIMEOUT (gosec G114)
func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")

		// SQL injection usage
		query := buildQuery(user)
		fmt.Println(query)

		// command injection
		runCommand(user)

		// 🔥 reflection XSS-like pattern
		fmt.Fprintf(w, "Hello %s", user)
	})

	http.ListenAndServe(":8080", nil) // G114
}

func main() {

	fmt.Println("Starting vulnerable app...")

	// 🔥 HARDCODED CREDENTIAL AGAIN
	dbPassword := "admin:admin123"
	fmt.Println(dbPassword)

	// 🔥 INSECURE RANDOM / WEAK LOGIC
	if password == "123" {
		fmt.Println("weak check")
	}

	startServer()
}
