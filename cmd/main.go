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
		t, err := template.ParseFiles("../static/template/index.html")
		if err != nil {
			log.Fatal("Error parsing template:", err)
		}
		t.Execute(w, nil)
	})

	// Untuk login, menggunakan http.HandleFunc
	r.Handle("/login", http.HandlerFunc(service.Login))
	r.Handle("/register", http.HandlerFunc(service.Register))

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
