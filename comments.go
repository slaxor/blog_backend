package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Handler function to get all comments for a blog post by ID
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Retrieve all comments for the post from the database
	rows, err := db.Query("SELECT * FROM comments WHERE post_id=?", postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Abrufen der Kommentare",
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	// Iterate over each row and add it to the response
	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Author,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Fehler beim Lesen des Kommentars",
				Data:    nil,
			})
			return
		}
		comments = append(comments, comment)
	}

	// Return the list of comments in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    comments,
	})
}

// Handler function to create a new comment for a blog post by ID
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	// Extract the comment data from the request body
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Ungültige Anforderung",
			Data:    nil,
		})
		return
	}

	// Insert the new comment into the database
	result, err := db.Exec("INSERT INTO comments (post_id, author, content) VALUES (?, ?, ?)", postID, comment.Author, comment.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Erstellen des Kommentars",
			Data:    nil,
		})
		return
	}

	// Get the ID of the newly created comment
	commentID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Fehler beim Abrufen der ID des Kommentars",
			Data:    nil,
		})
		return
	}
	comment.ID = int(commentID)
	comment.PostID = postID

	// Return the new comment in the response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "",
		Data:    comment,
	})
}
