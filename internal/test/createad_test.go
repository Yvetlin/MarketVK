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

func TestCreateAdHandler(t *testing.T) {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&model.User{}, &model.Ad{})
	dbpkg.DB = testDB
	router := handler.NewRouter()

	reqBody, _ := json.Marshal(map[string]string{
		"login":    "TestUser3",
		"password": "TestPassword3",
	})

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	loginBody, _ := json.Marshal(map[string]string{
		"login":    "TestUser3",
		"password": "TestPassword3",
	})
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["token"]

	adBody, _ := json.Marshal(map[string]interface{}{
		"title":       "Test Ad",
		"description": "Testing",
		"image_url":   "https://test.png",
		"price":       999.99,
	})

	req = httptest.NewRequest("POST", "/ads", bytes.NewReader(adBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected %v but got %v", http.StatusOK, w.Code)
	}

	var adResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &adResp)
	if adResp["ID"] == nil {
		t.Fatalf("No ID in response")
	}
}
