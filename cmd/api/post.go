package main

import (
	"log"
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"social-api/model"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	utils.WriteJson(w, http.StatusCreated, payload)
}

func (app *Application) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParams := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	ctx := r.Context()
	post, err := app.Store.Posts.GetByID(ctx, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			utils.WriteJsonError(w, http.StatusNotFound, "Post not found")
		default:
			log.Println("Error retrieving post:", err)
			utils.WriteJsonError(w, http.StatusInternalServerError, "Failed to retrieve post")
		}
		return
	}
	utils.WriteJson(w, http.StatusOK, post)
}
