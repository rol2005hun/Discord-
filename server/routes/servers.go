package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svelte/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateServerRequest definálja a szerver létrehozásához szükséges adatokat
type CreateServerRequest struct {
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

// CreateServerHandler kezeli a szerver létrehozását
func CreateServerHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Érvénytelen metódus. (POST szükséges)", http.StatusUnauthorized)
		return
	}

	var req CreateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Hiba a JSON dekódolásakor: %v", err), http.StatusBadRequest)
		return
	}

	// JWT tokenből megszerezzük a felhasználó ID-jét
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token hiányzik.", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Érvénytelen token.", http.StatusUnauthorized)
		return
	}

	// OwnerID-t állítjuk a tokenből jött felhasználói ID-re
	ownerID, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		http.Error(w, "Érvénytelen Owner ID.", http.StatusBadRequest)
		return
	}

	// Új Server típusú változó létrehozása az adatokkal
	server := models.Server{
		Name:      req.Name,
		OwnerID:   ownerID, // Beállítjuk az OwnerID-t a felhasználó azonosítójával
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Új szerver létrehozása
	result, err := models.CreateServer(server, dbClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Visszaadjuk a létrehozott szerver adatokat
	json.NewEncoder(w).Encode(result)
}

// GetServersHandler kezeli az összes szerver lekérdezését
func GetServersHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Érvénytelen metódus. (GET szükséges)", http.StatusUnauthorized)
		return
	}

	servers, err := models.GetServers(dbClient)
	if err != nil {
		http.Error(w, "Hiba az adatok lekérdezésében.", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(servers)
}

func GetUserServersHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Érvénytelen metódus. (GET szükséges)", http.StatusUnauthorized)
		return
	}

	// JWT tokenből megszerezzük a felhasználó ID-jét
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token hiányzik.", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Érvénytelen token.", http.StatusUnauthorized)
		return
	}

	// OwnerID-t állítjuk a tokenből jött felhasználói ID-re
	ownerID, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		http.Error(w, "Érvénytelen Owner ID.", http.StatusBadRequest)
		return
	}

	// Felhasználóhoz tartozó szerverek lekérdezése
	servers, err := models.GetServersByOwnerID(ownerID, dbClient)
	if err != nil {
		http.Error(w, "Hiba az adatok lekérdezésében.", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(servers)
}

// UpdateServerHandler kezeli a szerver adatok frissítését
func UpdateServerHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodPut {
		http.Error(w, "Érvénytelen metódus. (PUT szükséges)", http.StatusUnauthorized)
		return
	}

	// Szerver frissítési logika (használhatsz dinamikus paramétereket az URL-ből pl. szerver ID alapján)
	// A részleteket a kérésből (request) kell kinyerni
}

// DeleteServerHandler kezeli egy szerver törlését
func DeleteServerHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Érvénytelen metódus. (DELETE szükséges)", http.StatusUnauthorized)
		return
	}

	// Szerver törlési logika
	// Az URL-ben levő szerver ID alapján dolgozol
}
