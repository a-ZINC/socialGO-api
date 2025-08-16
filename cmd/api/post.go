package main

import (	
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"social-api/model"
	"strconv"
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
type PostWithMetadata struct {
	Posts  []*model.PostWithMetadata `json:"posts"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
	Order  string                    `json:"order"`
	Search string                    `json:"search"`
	Since  *string                   `json:"since"`
	Tags   []string                  `json:"tags"`
	Count  int64                     `json:"count"`
}


// CreatePosthandler creates a new post
// @Summary Create a new post
// @Description Create a new post with title, content, and tags
// @Tags Posts
// @Accept json
// @Produce json
// @Param body body PostPayload true "Post data"
// @Success 201 {object} PostPayload
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/posts [post]
func (app *Application) CreatePosthandler(w http.ResponseWriter, r *http.Request) {
	payload := &PostPayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		app.Logger.Errorf("Error reading post payload: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}

	app.Logger.Infof("Creating post with title: %v", payload)
	if err := utils.Validator.Struct(payload); err != nil {
		app.Logger.Errorf("Validation error: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err := app.Store.Posts.Create(ctx, &model.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	})
	if err != nil {
		app.Logger.Errorf("Error creating post: %v", err)
		app.Err.InternalServerError(w, r, err)
		return
	}
	app.Logger.Infof("Post created successfully: %v", payload)
	utils.JsonResponse(w, http.StatusCreated, payload)
}

// GetPostByIDHandler retrieves a post by its ID
// @Summary Get post by ID
// @Description Get post details by post ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Success 200 {object} model.PostWithMetadata
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/posts/{postId} [get]
func (app *Application) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)
	app.Logger.Infof("Retrieving post with ID: %d", id)
	ctx := r.Context()
	post, err := app.Store.Posts.GetByID(ctx, id)
	app.Logger.Infof("Post retrieved: %v", post)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Logger.Errorf("Post not found: %v", err)
			app.Err.NotFoundError(w, r, err)
		default:
			app.Logger.Errorf("Error retrieving post: %v", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	comments, _ := app.Store.Comments.GetByPostId(ctx, id)
	app.Logger.Infof("Comments retrieved for post ID %d: %v", id, comments)
	post.Comments = comments
	app.Logger.Infof("Returning post with ID %d: %v", id, post)
	utils.WriteJson(w, http.StatusOK, post)
}

func (app *Application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)
	app.Logger.Infof("Deleting post with ID: %d", id)

	ctx := r.Context()
	err := app.Store.Posts.Delete(ctx, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Logger.Errorf("Post not found for deletion: %v", err)
			app.Err.NotFoundError(w, r, err)
		default:
			app.Logger.Errorf("Error retrieving post: %v", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	app.Logger.Infof("Post with ID %d deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}


// UpdatePostHandler updates an existing post
// @Summary Update an existing post
// @Description Update a post by its ID with new title and/or content
// @Tags Posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param body body UpdatePayload true "Post data"
// @Success 200 {object} UpdatePayload
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/posts/{postId} [patch]
func (app *Application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	id := app.Middleware.GetPostIdFromContext(r)
	app.Logger.Infof("Updating post with ID: %d", id)
	post, err := app.Store.Posts.GetByID(r.Context(), id)
	app.Logger.Infof("Post retrieved for update: %v", post)
	if err != nil {
		app.Logger.Errorf("Error retrieving post for update: %v", err)
		app.Err.NotFoundError(w, r, err)
		return
	}

	payload := &UpdatePayload{}
	if err := utils.ReadJson(w, r, payload); err != nil {
		app.Logger.Errorf("Error reading update payload: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	app.Logger.Infof("Updating post with ID %d with payload: %v", id, payload)

	if err := utils.Validator.Struct(payload); err != nil {
		app.Logger.Errorf("Validation error: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	app.Logger.Infof("Post content after update: %s", post.Content)
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	app.Logger.Infof("Post after update: %v", post)

	ctx := r.Context()
	err = app.Store.Posts.Update(ctx, post, id)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Logger.Errorf("Post not found for update: %v", err)
			app.Err.NotFoundError(w, r, err)
		default:
			app.Logger.Errorf("Error updating post: %v", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	app.Logger.Infof("Post with ID %d updated successfully", id)
	if err := utils.WriteJson(w, http.StatusOK, payload); err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}
}


// GetUserFeedHandler retrieves the user's feed
// @Summary Get user feed
// @Description Get the feed of posts for a user
// @Tags Users
// @Accept json
// @Produce json
// @Param pagination query store.Pagination true "Pagination parameters"
// @Param userId path int true "User ID"
// @Success 200 {object} PostWithMetadata
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/user/feed [get]
func (app *Application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pagination := &store.Pagination{
		Limit:  10,
		Offset: 0,
		Order:  "desc",
		Search: "",
		Since:  nil,
		Tags:   []string{},
	}
	app.Logger.Infof("Retrieving user feed with pagination: %v", pagination)
	pagination, err := pagination.GetPaginated(r)
	if err != nil {
		app.Logger.Errorf("Error retrieving pagination: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	app.Logger.Infof("Pagination after retrieval: %v", pagination)

	if err := utils.Validator.Struct(pagination); err != nil {
		app.Logger.Errorf("Validation error for pagination: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	app.Logger.Infof("Pagination after validation: %v", pagination)

	ctx := r.Context()
	userId, err := strconv.ParseInt("10", 10, 64)
	if err != nil {
		app.Logger.Errorf("Failed to parse user ID: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	cnt, err := app.Store.Posts.GetUserFeedCount(ctx, userId, pagination)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Logger.Errorf("User feed not found: %v", err)
			app.Err.NotFoundError(w, r, err)
		default:
			app.Logger.Errorf("Error retrieving user feed count: %v", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	app.Logger.Infof("User feed count retrieved: %d", cnt)
	posts, err := app.Store.Posts.GetUserFeed(ctx, userId, pagination)
	app.Logger.Infof("Posts retrieved for user feed: %v", posts)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Logger.Errorf("No posts found for user feed: %v", err)
			app.Err.NotFoundError(w, r, err)
		default:
			app.Logger.Errorf("Error retrieving post: %v", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	postWithMetadata := PostWithMetadata{
		Posts:  posts,
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Order:  pagination.Order,
		Search: pagination.Search,
		Since:  pagination.Since,
		Tags:   pagination.Tags,
		Count:  cnt,
	}
	app.Logger.Infof("Returning user feed with metadata: %v", postWithMetadata)
	if err := utils.WriteJson(w, http.StatusOK, postWithMetadata); err != nil {
		app.Logger.Errorf("Error writing user feed response: %v", err)
		app.Err.InternalServerError(w, r, err)
		return
	}
}
