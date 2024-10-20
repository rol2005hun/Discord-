package main

import (
	"fmt"
	"log"
	"net/http"
	"svelte/config"
	"svelte/middlewares"
	"svelte/routes"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Adatbázis kapcsolat inicializálása
	dbClient := config.ConnectDB()
	defer config.DisconnectDB(dbClient)

	mux := http.NewServeMux()

	// users
	mux.HandleFunc("/users/getAll", func(w http.ResponseWriter, r *http.Request) {
		routes.UsersHandler(w, r, dbClient)
	})

	mux.Handle("/users/createUser", middlewares.UserValidationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.CreateUserHandler(w, r, dbClient)
	}), dbClient))

	mux.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		routes.LoginHandler(w, r, dbClient)
	})

	mux.Handle("/users/logout", middlewares.TokenValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.LogoutHandler(w, r)
	})))

	// servers
	mux.HandleFunc("/servers/createServer", func(w http.ResponseWriter, r *http.Request) {
		routes.CreateServerHandler(w, r, dbClient)
	})

	mux.HandleFunc("/servers/getAll", func(w http.ResponseWriter, r *http.Request) {
		routes.GetServersHandler(w, r, dbClient)
	})

	// mux.HandleFunc("/servers/getServer", func(w http.ResponseWriter, r *http.Request) {
	// 	routes.GetServerHandler(w, r, dbClient)
	// })

	mux.HandleFunc("/servers/updateServer", func(w http.ResponseWriter, r *http.Request) {
		routes.UpdateServerHandler(w, r, dbClient)
	})

	// mux.HandleFunc("/servers/deleteServer", func(w http.ResponseWriter, r *http.Request) {
	// 	routes.DeleteServerHandler(w, r, dbClient)
	// })

	mux.HandleFunc("/servers/getUserServers", func(w http.ResponseWriter, r *http.Request) {
		routes.GetUserServersHandler(w, r, dbClient)
	})

	handler := CORS(mux)

	fmt.Println("[Wolimby] A szerver elindult a 8080-as porton.")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
