package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}
	// Benutzername und Passwort aus dem Anfragekörper extrahieren
	// username := r.FormValue("username")
	// password := r.FormValue("password")

	// Benutzer in der SQLite-Datenbank suchen
	// var user User
	err = db.QueryRow("SELECT id, name, password_hash FROM users WHERE name=?", user.Name).Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültiger Benutzername oder Passwort",
			Data:    nil,
		})
		return
	}

	// Passwort überprüfen
	if !user.Authenticate("mypassword") {
		// incorrect password
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültiger Benutzername oder Passwort",
			Data:    nil,
		})
		return
	}

	// JWT-Token generieren
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("geheimes_token"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Generieren des Tokens",
			Data:    nil,
		})
		return
	}

	// JWT-Token an den Benutzer zurückgeben
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data: map[string]string{
			"token": tokenString,
		},
	})
}
