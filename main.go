package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// SQLite Datenbankverbindung herstellen
	db, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// HTTP-Handler einrichten
	router := mux.NewRouter()
	router.HandleFunc("/api/posts", GetPostsHandler).Methods("GET")
	router.HandleFunc("/api/posts", AuthMiddleware(CreatePostHandler)).Methods("POST")
	router.HandleFunc("/api/posts/{id:[0-9]+}", GetPostHandler).Methods("GET")
	router.HandleFunc("/api/posts/{id:[0-9]+}", AuthMiddleware(UpdatePostHandler)).Methods("PUT")
	router.HandleFunc("/api/posts/{id:[0-9]+}", AuthMiddleware(DeletePostHandler)).Methods("DELETE")
	router.HandleFunc("/api/posts/{id:[0-9]+}/comments", GetCommentsHandler).Methods("GET")
	router.HandleFunc("/api/posts/{id:[0-9]+}/comments", AuthMiddleware(CreateCommentHandler)).Methods("POST")
	router.HandleFunc("/api/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/token", GenerateJWTHandler).Methods("POST")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("geheimes_token"))

	log.Print(tokenString)

	logger := handlers.LoggingHandler(os.Stdout, router)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	logger = handlers.CORS(headers, methods, origins)(logger)
	server := &http.Server{
		Addr:         "127.0.0.1:18000",
		Handler:      logger,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Starting server on %s...", server.Addr)
	log.Fatal(server.ListenAndServe())
	// http.ListenAndServe("", handlers.CORS(headers, methods, origins)(router))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWT-Token aus dem Autorisierungs-Header extrahieren
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Kein JWT-Token angegeben",
				Data:    nil,
			})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// JWT-Token überprüfen
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Ungültiger Signaturalgorithmus")
			}
			return []byte("geheimes_token"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		// JWT-Claims überprüfen
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(float64)
			r = r.WithContext(context.WithValue(r.Context(), "user_id", int64(userID)))
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Ungültiger JWT-Token",
				Data:    nil,
			})
			return
		}
	}
}
