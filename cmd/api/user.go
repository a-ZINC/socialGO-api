package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"social-api/model"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FollowUser struct {
	UserId int64 `json:"user_id"`
}
type UserContextKey string

var key UserContextKey = "user"

// GetUserByIDHandler retrieves a user by their ID from the context
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/user/{userId} [get]
func (app *Application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.GetUserFromContext(r)
	app.Logger.Infof("Retrieving user with ID: %d", user.ID)
	if err != nil {
		app.Logger.Errorf("Error retrieving user from context: %v", err)
		app.Err.NotFoundError(w, r, err)
		return
	}
	app.Logger.Infof("Retrieving user with ID: %d", user.ID)
	utils.WriteJson(w, http.StatusOK, user)
}

// FollowUserHandler follows a user
// @Summary Follow a user
// @Description Follow a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param body body FollowUser true "User to follow"
// @Success 200 {object} map[string]string
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/user/{userId}/follow [put]
func (app *Application) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, err := app.GetUserFromContext(r)
	app.Logger.Infof("Follower user ID: %d", followerUser.ID)
	if err != nil {
		app.Logger.Errorf("Error retrieving follower user from context: %v", err)
		app.Err.NotFoundError(w, r, err)
		return
	}
	var followUser FollowUser
	if err := utils.ReadJson(w, r, &followUser); err != nil {
		app.Logger.Errorf("Error reading follow user payload: %v", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	followerUserId := followerUser.ID
	app.Logger.Infof("Follower user ID: %d, Follow user ID: %d", followerUserId, followUser.UserId)

	err = app.Store.Follower.Follow(r.Context(), followUser.UserId, followerUserId)
	if err != nil {
		app.Logger.Errorf("Error following user: %v", err)
		app.Err.InternalServerError(w, r, err)
		return
	}
	app.Logger.Infof("User with ID %d followed successfully", followUser.UserId)

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Followed successfully"})
}

// UnfollowUserHandler unfollows a user
// @Summary Unfollow a user
// @Description Unfollow a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param body body FollowUser true "User to unfollow"
// @Success 200 {object} map[string]string
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/user/{userId}/unfollow [put]
func (app *Application) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowerUser, err := app.GetUserFromContext(r)
	if err != nil {
		app.Err.NotFoundError(w, r, err)
		return
	}

	var unfollowUser FollowUser
	if err := utils.ReadJson(w, r, &unfollowUser); err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}
	unfollowerUserId := unfollowerUser.ID
	err = app.Store.Follower.Unfollow(r.Context(), unfollowerUserId, unfollowUser.UserId)
	if err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Unfollowed successfully"})
}
func (app *Application) CreateUserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.Err.BadRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.Store.Users.GetByID(ctx, userId)
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
		ctx = context.WithValue(ctx, key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *Application) GetUserFromContext(r *http.Request) (model.User, error) {
	user, ok := r.Context().Value(key).(model.User)
	if !ok {
		return model.User{}, errors.New("user not found in context")
	}
	return user, nil
}

// Activate User
// @Summary Activate a user
// @Description Activate a user by their token
// @Tags Users
// @Accept json
// @Produce json
// @Param token path string true "Activation token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /v1/user/activate/{token} [put]
func (app *Application) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Activating user")
	token := chi.URLParam(r, "token")
	log.Printf("Activating user with token: %s", token)
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	log.Printf("Activating user with token: %s", hashedToken)
	err := app.Store.Users.ActivateUser(ctx, hashedToken)
	if err != nil {
		switch err {
		case store.ErrPostNotFound:
			app.Err.NotFoundError(w, r, err)
		default:
			app.Err.InternalServerError(w, r, err)
		}
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "User activated successfully"})
}
