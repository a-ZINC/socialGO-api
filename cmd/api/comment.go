package main

import (
	"net/http"
	"social-api/cmd/utils"
	"social-api/model"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
}

// CreateCommentHandler creates a new comment on a post
// @Summary Create a new comment
// @Description Create a new comment on a post by post ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param body body CommentPayload true "Comment data"
// @Success 201 {object} CommentPayload
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/posts/{postId}/comments [post]
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
		UserID:  1,
		PostID:  postId,
		Content: payload.Content,
	})

	if err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, payload)
}
