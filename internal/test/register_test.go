package test

import (
	"bytes"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	dbpkg "MarketVK/internal/db"
	"MarketVK/internal/handler"
)

func TestRegisterHandler(t *testing.T) {

	if err := dbpkg.InitTestDB(); err != nil {
		t.Fatal(err)
	}

	router := handler.NewRouter()

	reqBody, _ := json.Marshal(map[string]string{
		"login":    "TestUser",
		"password": "TestPassword",
	})

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d\n", http.StatusOK, w.Code)
	}
}
