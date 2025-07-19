package handler

import (
	"encoding/json"
	"net/http"

	"strconv"
	"strings"

	"MarketVK/internal/db"
	"MarketVK/internal/model"
	"MarketVK/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPassLength  = 6
	MinLoginLength = 4
)

var (
	users  = map[string]model.User{}
	userID = 1
	ads    = []model.Ad{}
	adID   = 1
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/ads", AdsHandler)
	return mux
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "login/password too short", http.StatusBadRequest)
		return
	}

	if len(req.Login) < MinLoginLength || len(req.Password) < MinPassLength {
		http.Error(w, "login/password too short", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	user := model.User{
		Login:    req.Login,
		Password: string(hashed),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "User exists or DB error", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var user model.User
	err := db.DB.Where("login = ?", req.Login).First(&user).Error
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func AdsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createAdHandler(w, r)
	case http.MethodGet:
		listAdsHandler(w, r)
	case http.MethodPut:
		UpdateAdHandler(w, r)
	case http.MethodDelete:
		DeleteAdHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createAdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	userID, err := auth.ParseToken(parts[1])
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		Price       float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ad := model.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		AuthorID:    userID,
	}
	if err := db.DB.Create(&ad).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}

func listAdsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	sortBy := r.URL.Query().Get("sort_by") // "date" или "price"
	order := r.URL.Query().Get("order")    // "asc" или "desc"

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := db.DB.Model(&model.Ad{})

	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("price >= ?", minPrice)
		}
	}
	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("price <= ?", maxPrice)
		}
	}

	if sortBy == "price" {
		if order == "asc" {
			query = query.Order("price asc")
		} else {
			query = query.Order("price desc")
		}
	} else {
		if order == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	var ads []model.Ad
	if err := query.Preload("Author").
		Offset(offset).
		Limit(limit).
		Find(&ads).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

func UpdateAdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ParseToken(parts[1])
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID          uint    `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var ad model.Ad
	if err := db.DB.First(&ad, req.ID).Error; err != nil {
		http.Error(w, "ad not found", http.StatusBadRequest)
		return
	}

	if ad.AuthorID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	ad.Title = req.Title
	ad.Description = req.Description
	ad.ImageURL = req.ImageURL
	ad.Price = req.Price

	if err := db.DB.Save(&ad).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}

func DeleteAdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ParseToken(parts[1])
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID uint `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var ad model.Ad
	if err := db.DB.First(&ad, req.ID).Error; err != nil {
		http.Error(w, "ad not found", http.StatusBadRequest)
		return
	}

	if ad.AuthorID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := db.DB.Delete(&ad).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
