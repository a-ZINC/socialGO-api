package main

import (
	"net/http"
	"social-api/cmd/utils"
	"social-api/model"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentPayload struct {
	Content string  `json:"content" validate:"required,max=100"`
}
func (app *Application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postIdParams := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(postIdParams, 10, 64)
	if err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}
	payload := &CommentPayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}

	err = utils.Validator.Struct(payload)
	if err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.Store.Comments.Create(ctx, &model.Comment{
		UserID: 1,
		PostID: postId,
		Content: payload.Content,
	})

	if err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, payload)
}