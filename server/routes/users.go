package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"svelte/logs"
	"svelte/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("beresildiko")

type Credentials struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// LoginHandler kezeli a bejelentkezést és a JWT generálást
func LoginHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Érvénytelen metódus. (POST szükséges)", http.StatusUnauthorized)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Hibás JSON.", http.StatusUnauthorized)
		return
	}

	// Ellenőrizzük a felhasználó létezését az adatbázisban
	user, err := models.GetUserByUsernameOrEmail(creds.UsernameOrEmail, dbClient)
	if err != nil || user.Password != creds.Password {
		http.Error(w, "Helytelen email vagy jelszó.", http.StatusUnauthorized)
		return
	}

	// JWT token generálás
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		ID:       user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Nem sikerült a token generálása.", http.StatusUnauthorized)
		return
	}

	// Token visszaküldése a felhasználónak
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  expirationTime,
	})

	json.NewEncoder(w).Encode(Response{
		Message: "Sikeres bejelentkezés.",
		Status:  http.StatusOK,
	})

	message := fmt.Sprintf("[%s] [Wolimby - Bejelentkezés] ID: %s, felhasználónév: %s", time.Now().Format(time.RFC3339), user.ID.Hex(), creds.UsernameOrEmail)
	logFilePath := "logs/logins.txt"

	if err := logs.Log(message, logFilePath); err != nil {
		log.Printf("Sikertelen bejelentkezés logolás: %v", err)
	}
}

// CreateUserHandler kezeli a POST kéréseket
func CreateUserHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Érvénytelen metódus. (POST szükséges)", http.StatusUnauthorized)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Hiba a JSON dekódolásakor: %v", err), http.StatusBadRequest)
		return
	}

	recordedUser := models.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := models.CreateUser(recordedUser, dbClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)

	message := fmt.Sprintf("[%s] [Wolimby - Regisztráció] Felhasználónév: %s, email: %s", time.Now().Format(time.RFC3339), req.Username, req.Email)
	logFilePath := "logs/registrations.txt"

	if err := logs.Log(message, logFilePath); err != nil {
		log.Printf("Sikertelen regisztráció logolás: %v", err)
	}
}

// UsersHandler kezeli a GET kéréseket
func UsersHandler(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	users, err := models.GetUsers(dbClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	if r.Method != http.MethodPost {
		http.Error(w, "Érvénytelen metódus. (POST szükséges)", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	message := fmt.Sprintf("[%s] [Wolimby - Kijelentkezés] ID: %s", time.Now().Format(time.RFC3339), userID)
	logFilePath := "logs/logouts.txt"

	if err := logs.Log(message, logFilePath); err != nil {
		log.Printf("Sikertelen kijelentkezés logolás: %v", err)
	}

	json.NewEncoder(w).Encode(Response{
		Message: "Sikeres kijelentkezés.",
		Status:  http.StatusOK,
	})
}
