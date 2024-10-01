package main

import (
	"fmt"
	"log"
	"net/http"

	"blog-api/internal/handler"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

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
