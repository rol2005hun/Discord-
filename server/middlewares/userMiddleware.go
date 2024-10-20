package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"svelte/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("beresildiko")

// ValidationError struct hibaüzenetekhez
type ValidationError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// Response struct válaszokhoz
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Claims struktúra
type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// User struct a bejövő felhasználói adatok tárolására
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWTMiddleware ellenőrzi a JWT érvényességét
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Nincs token.", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Hiba a token olvasásakor.", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		// Token érvényességének ellenőrzése
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Érvénytelen aláírás.", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Hibás token.", http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			http.Error(w, "Érvénytelen token.", http.StatusUnauthorized)
			return
		}

		// Token lejáratának ellenőrzése
		if claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "Lejárt token.", http.StatusUnauthorized)
			return
		}

		// Továbbítás a következő handlerhez
		next.ServeHTTP(w, r)
	})
}

// UserValidationMiddleware ellenőrzi a felhasználói adatokat
func UserValidationMiddleware(next http.Handler, dbClient *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Olvassa be a kérés törzsét
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			// JSON adatok deszerializálása
			var user User
			if err := json.Unmarshal(body, &user); err != nil {
				http.Error(w, "Hibás JSON.", http.StatusBadRequest)
				return
			}

			// Hosszúság ellenőrzések
			if len(user.Username) > 15 || len(user.Username) < 3 {
				http.Error(w, "A felhasználónév 3 és 15 karakter között kell legyen.", http.StatusBadRequest)
				return
			}

			if len(user.Name) > 50 {
				http.Error(w, "A név maximum 50 karakter lehet.", http.StatusBadRequest)
				return
			}

			if len(user.Password) < 8 {
				http.Error(w, "A jelszó legalább 8 karakter hosszú kell legyen.", http.StatusBadRequest)
				return
			}

			// Felhasználónév és email ellenőrzése
			usernameResult, err := models.GetUserByUsernameOrEmail(user.Username, dbClient)
			if err == nil {
				if usernameResult.Username == user.Username {
					http.Error(w, "A felhasználónév már foglalt.", http.StatusBadRequest)
					return
				}
			}

			emailResult, err := models.GetUserByUsernameOrEmail(user.Email, dbClient)
			if err == nil {
				if emailResult.Email == user.Email {
					http.Error(w, "Az email cím már foglalt.", http.StatusBadRequest)
					return
				}
			}

			// RegEx az email validálásához
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			isValidEmail := regexp.MustCompile(emailRegex).MatchString(user.Email)

			if !isValidEmail {
				http.Error(w, "Érvénytelen email cím.", http.StatusBadRequest)
				return
			}

			// Állítsa vissza a törzset, hogy a következő handler olvashassa
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r) // Hívja a következő handlert
	})
}

func TokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
