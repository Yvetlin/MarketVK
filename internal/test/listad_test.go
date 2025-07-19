package test

import (
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

func TestListAdsHandler(t *testing.T) {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&model.User{}, &model.Ad{})
	dbpkg.DB = testDB
	router := handler.NewRouter()

	user := model.User{Login: "u", Password: "p"}
	testDB.Create(&user)
	ad := model.Ad{Title: "T", Description: "D", ImageURL: "url", Price: 1, AuthorID: int(user.ID)}
	testDB.Create(&ad)

	req := httptest.NewRequest("GET", "/ads", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var ads []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &ads)
	if len(ads) == 0 {
		t.Fatalf("No ads found")
	}
}
