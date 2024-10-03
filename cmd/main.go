package main

import (
	"blog-api/internal/handler"
	"blog-api/internal/service"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		t, err := template.ParseFiles("../static/template/index.html")
		if err != nil {
			log.Fatal("Error parsing template:", err)
		}
		t.Execute(w, nil)
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
	r.Handle("/posts", service.AuthMiddleware(http.HandlerFunc(handler.GetPost))).Methods("GET")
	r.Handle("/posts/{id}", service.AuthMiddleware(http.HandlerFunc(handler.GetPostbyID))).Methods("GET")
	r.Handle("/posts/create", service.AuthMiddleware(http.HandlerFunc(handler.CreatePost))).Methods("POST")
	r.Handle("/posts/update/{id}", service.AuthMiddleware(http.HandlerFunc(handler.UpdatePost))).Methods("PUT")
	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
