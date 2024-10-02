package main

import (
	"blog-api/internal/handler"
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
	r.HandleFunc("/posts", handler.GetPost).Methods("GET")
	r.HandleFunc("/posts/{id}", handler.GetPostbyID).Methods("GET")
	r.HandleFunc("/posts/create", handler.CreatePost).Methods("POST")
	r.HandleFunc("/posts/update/{id}", handler.UpdatePost).Methods("PUT")

	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
