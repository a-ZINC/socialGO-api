package main

import (
	"log"
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"social-api/model"
)

type PostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}
type UpdatePayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *Application) CreatePosthandler(w http.ResponseWriter, r *http.Request) {
	payload := &PostPayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		log.Println("Creating post:", payload)
		app.Err.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		app.Err.BadRequestError(w, r, err)
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
		app.Err.InternalServerError(w, r, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, payload)
}

func (app *Application) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)
	ctx := r.Context()
	post, err := app.Store.Posts.GetByID(ctx, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Err.NotFoundError(w, r, err)
		default:
			log.Println("Error retrieving post:", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	comments, _ := app.Store.Comments.GetByPostId(ctx, id)
	post.Comments = comments
	utils.WriteJson(w, http.StatusOK, post)
}

func (app *Application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)

	ctx := r.Context()
	err := app.Store.Posts.Delete(ctx, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Err.NotFoundError(w, r, err)
		default:
			log.Println("Error retrieving post:", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)
	post, err := app.Store.Posts.GetByID(r.Context(), id)
	if err != nil {
		app.Err.NotFoundError(w, r, err)
		return
	}

	payload := &UpdatePayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	ctx := r.Context()
	err = app.Store.Posts.Update(ctx, post, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Err.NotFoundError(w, r, err)
		default:
			log.Println("Error retrieving post:", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, payload); err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}
}
