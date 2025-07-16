package handler

import (
	"MarketVK/internal/model"
	"MarketVK/pkg/auth"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	if _, exists := users[req.Login]; exists {
		http.Error(w, "user exists", http.StatusConflict)
		return
	}

	user := model.User{
		ID:       userID,
		Login:    req.Login,
		Password: req.Password,
	}
	users[req.Login] = user
	userID++

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

	user, exists := users[req.Login]
	if !exists || user.Password != req.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
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
		listAdsHandler(w, r) //!!
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createAdHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(parts[1])
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

	if len(req.Title) == 0 || len(req.Title) > 100 || req.Price < 0 {
		http.Error(w, "invalid ad params", http.StatusBadRequest)
		return
	}

	ad := model.Ad{
		ID:          adID,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		AuthorID:    userID,
	}
	adID++
	ads = append(ads, ad)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}

func listAdsHandler(w http.ResponseWriter, r *http.Request) {
	// Можно реализовать пагинацию через query: ?page=1&limit=10
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}
