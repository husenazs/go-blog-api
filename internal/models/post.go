package models

type Post struct {
	Post_ID    string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"body"`
	Author     string `json:"author"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
