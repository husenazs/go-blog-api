package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"body"`
	Author  string `json:"author"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/husenDB")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database connection is not alive:", err)
	}
}

var responses []map[string]interface{}

func createPost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var newPost Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert post into database
	stmt, err := db.Prepare("INSERT INTO POSTS (TITLE, CONTENT, AUTHOR) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newPost.Title, newPost.Content, newPost.Author)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newPost.ID = int(id)
	w.WriteHeader(http.StatusCreated)

	// Tambahkan response ke array responses
	responses = append(responses, map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": http.StatusCreated,
		"data":   newPost,
	})
	json.NewEncoder(w).Encode(responses)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT ID, TITLE, CONTENT FROM POSTS")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func getPostbyID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post Post
	err = db.QueryRow("SELECT ID, TITLE, CONTENT FROM POSTS WHERE ID = ?", id).Scan(&post.ID, &post.Title, &post.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tambahkan response ke array responses
	responses = append(responses, map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": http.StatusCreated,
		"data":   post,
	})

	json.NewEncoder(w).Encode(responses)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE POSTS SET TITLE = ?, CONTENT = ? WHERE ID = ?")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(post.Title, post.Content, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if affectedRows == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var arrReposnses []map[string]interface{}
	arrReposnses = append(arrReposnses, map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": http.StatusOK,
		"data":   post,
	})
	json.NewEncoder(w).Encode(arrReposnses)

}

func main() {
	http.HandleFunc("/posts", getPost)
	http.HandleFunc("/posts/", getPostbyID)
	http.HandleFunc("/posts/create", createPost)
	http.HandleFunc("/posts/update/", updatePost)

	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
