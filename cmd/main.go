package main

import (
	"blog-api/internal/handler"
	"blog-api/internal/models"
	"blog-api/internal/service"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func formattedTime(timeStr string) string {
	// Format datetime sesuai kebutuhan
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Println("Error parsing time:", err)
		return timeStr
	}
	formattedTime := fmt.Sprintf("%d-%s-%d", t.Day(), t.Month().String(), t.Year())

	return formattedTime
}

func randNumber() string {
	return fmt.Sprintf("%d", time.Now().Day())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Call API untuk mendapatkan posts
		resp, err := http.Get("http://localhost:8080/posts")
		if err != nil {
			log.Println("Error getting posts:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var posts []models.Post
		// Decode response body into slice of posts
		if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
			log.Println("Error decoding posts:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Loop through the posts and update the Created_at and Updated_at fields
		for i, post := range posts {
			posts[i].Created_at = formattedTime(post.Created_at)
			posts[i].Updated_at = formattedTime(post.Updated_at)
			posts[i].Post_ID = randNumber()
		}

		// Parsing template dan passing data posts
		t, err := template.ParseFiles("../static/template/index.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Eksekusi template dengan data posts
		if err := t.Execute(w, struct {
			Posts []models.Post
		}{Posts: posts}); err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Cek cookie, jika ada token JWT maka redirect ke /
		_, err := r.Cookie("token")
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Jika tidak ada cookie, tampilkan halaman login
		t, err := template.ParseFiles("../static/template/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error parsing template:", err)
			return
		}
		t.Execute(w, nil)
	}).Methods("GET")
	r.Handle("/logout", http.HandlerFunc(service.Logout)).Methods("GET")

	// Untuk login, menggunakan http.HandleFunc
	r.Handle("/login", http.HandlerFunc(service.Login)).Methods("POST")
	r.Handle("/register", http.HandlerFunc(service.Register)).Methods("POST")

	// Menggunakan http.HandlerFunc untuk mengkonversi handler ke http.Handler
	r.Handle("/posts", http.HandlerFunc(handler.GetPost)).Methods("GET")
	r.Handle("/posts/{id}", http.HandlerFunc(handler.GetPostbyID)).Methods("GET")
	r.Handle("/posts/create", service.AuthMiddleware(http.HandlerFunc(handler.CreatePost))).Methods("POST")
	r.Handle("/posts/update/{id}", service.AuthMiddleware(http.HandlerFunc(handler.UpdatePost))).Methods("PUT")
	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
