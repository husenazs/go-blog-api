package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"blog-api/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

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

// Fungsi CreatePost untuk menambahkan post baru
func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	author := newPost.Author
	currentTime := time.Now() // Mengambil waktu sekali saja
	id := author + strconv.FormatInt(currentTime.UnixNano(), 10)
	newPost.Post_ID = id

	// Format datetime menjadi "YYYY-MM-DD HH:MM:SS"
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// Insert post into database
	stmt, err := db.Prepare("INSERT INTO POSTS (POST_ID, TITLE, CONTENT, AUTHOR, CREATED_AT, UPDATED_AT) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPost.Post_ID, newPost.Title, newPost.Content, newPost.Author, formattedTime, formattedTime)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newPost.Created_at = formattedTime
	newPost.Updated_at = formattedTime

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

// Fungsi GetPost untuk mendapatkan semua post
func GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT POST_ID, TITLE, CONTENT, AUTHOR, CREATED_AT, UPDATED_AT FROM POSTS")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.Post_ID, &post.Title, &post.Content, &post.Author, &post.Created_at, &post.Updated_at)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

// Fungsi GetPostbyID untuk mendapatkan post berdasarkan ID
func GetPostbyID(w http.ResponseWriter, r *http.Request) {
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

	var post models.Post
	err = db.QueryRow("SELECT * FROM POSTS WHERE ID = ?", id).Scan(&post.Post_ID, &post.Title, &post.Content, &post.Author, &post.Created_at, &post.Updated_at)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// Fungsi UpdatePost untuk memperbarui post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	var post models.Post
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

	json.NewEncoder(w).Encode(post)
}
