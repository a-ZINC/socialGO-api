package main

import (
	"log"
	"net/http"
	"social-api/cmd/utils"
	"social-api/model"
)

type PostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *Application) CreatePosthandler(w http.ResponseWriter, r *http.Request) {
	payload := &PostPayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		log.Println("Creating post:", payload)
		if writeErr := utils.WriteJsonError(w, http.StatusBadRequest, "Invalid request payload"); writeErr != nil {
			http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		}
		return 
	}
	ctx := r.Context()
	err := app.Store.Posts.Create(ctx, model.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	})
	if err != nil {
		if writeErr := utils.WriteJsonError(w, http.StatusInternalServerError, "Failed to create post"); writeErr != nil {
			http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		}
		return
	}
	utils.WriteJson(w, http.StatusCreated, payload);
}
