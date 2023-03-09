package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}
	var user User
	err = db.QueryRow("SELECT password_hash FROM users WHERE name=?", login.Username).Scan(&user.PasswordHash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültiger Benutzername oder Passwort",
			Data:    nil,
		})
		return
	}

	// log.Print(user.PasswordHash)
	// user.SetPassword("foobar")
	// log.Print(user.PasswordHash)
	if !user.Authenticate(login.Password) {
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

/*
hash, err := HashPassword("mypassword")
if err != nil {
    // handle error
}

_, err = db.Exec("INSERT INTO users (name, password_hash) VALUES (?, ?)", "slaxor", hash)
if err != nil {
    // handle error
}

var passwordHash string
err := db.QueryRow("SELECT password_hash FROM users WHERE name = ?", "slaxor").Scan(&passwordHash)
if err != nil {
    // handle error
}

if !CheckPasswordHash("mypassword", passwordHash) {
    // incorrect password
} else {
    // correct password
}

*/
