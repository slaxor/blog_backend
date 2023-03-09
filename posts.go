package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Handler function to get all blog posts
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all posts from the database
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Abrufen der Beiträge",
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	// Iterate over each row and add it to the response
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Fehler beim Lesen des Beitrags",
				Data:    nil,
			})
			return
		}
		posts = append(posts, post)
	}

	// Return the list of posts in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    posts,
	})
}

// Handler function to create a new blog post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post data from the request body
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}

	// Insert the new post into the database
	result, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", post.Title, post.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Erstellen des Beitrags",
			Data:    nil,
		})
		return
	}

	// Get the ID of the newly created post
	postID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Abrufen der ID des Beitrags",
			Data:    nil,
		})
		return
	}
	post.ID = int(postID)

	// Return the new post in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    post,
	})
}

// Handler function to get a single blog post by ID
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL parameter
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}

	// Retrieve the post from the database by ID
	var post Post
	err = db.QueryRow("SELECT * FROM posts WHERE id=?", postID).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Beitrag nicht gefunden",
				Data:    nil,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Abrufen des Beitrags",
			Data:    nil,
		})
		return
	}

	// Return the post in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    post,
	})

}

// Handler function to update an existing blog post
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL parameter
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}

	// Extract the updated post data from the request body
	var post Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}

	// Update the post in the database by ID
	_, err = db.Exec("UPDATE posts SET title=?, content=? WHERE id=?", post.Title, post.Content, postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Aktualisieren des Beitrags",
			Data:    nil,
		})
		return
	}

	// Return the updated post in the response
	post.ID = postID
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    post,
	})
}

// Handler function to delete an existing blog post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL parameter
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}
	// Delete the post from the database by ID
	_, err = db.Exec("DELETE FROM posts WHERE id=?", postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Löschen des Beitrags",
			Data:    nil,
		})
		return
	}

	// Return a success message in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Beitrag erfolgreich gelöscht",
		Data:    nil,
	})

	// Delete the post from the database by ID
	_, err = db.Exec("DELETE FROM posts WHERE id=?", postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Löschen des Beitrags",
			Data:    nil,
		})
		return
	}

	// Return a success message in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Beitrag erfolgreich gelöscht",
		Data:    nil,
	})
}
