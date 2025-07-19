package test

import (
	"bytes"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	dbpkg "MarketVK/internal/db"
	"MarketVK/internal/handler"
	"MarketVK/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLoginHandler(t *testing.T) {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&model.User{}, &model.Ad{})
	dbpkg.DB = testDB
	router := handler.NewRouter()

	reqBody, _ := json.Marshal(map[string]string{
		"login":    "TestUser2",
		"password": "TestPassword2",
	})

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal("Registration failed")
	}

	loginBody, _ := json.Marshal(map[string]string{
		"login":    "TestUser2",
		"password": "TestPassword2",
	})
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Fatal("No token returned")
	}
}
