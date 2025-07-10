package model

type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}