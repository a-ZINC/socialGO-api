package main

import (
	"context"
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

func (app *Application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.GetUserFromContext(r)
	if err != nil {
		app.Err.NotFoundError(w, r, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, user)
}

func (app *Application) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, err := app.GetUserFromContext(r)
	if err != nil {
		app.Err.NotFoundError(w, r, err)
		return
	}
	var followUser FollowUser
	if err := utils.ReadJson(w, r, &followUser); err != nil {
		app.Err.BadRequestError(w, r, err)
		return
	}
	followerUserId := followerUser.ID

	err = app.Store.Follower.Follow(r.Context(), followerUserId, followUser.UserId)
	if err != nil {
		app.Err.InternalServerError(w, r, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Followed successfully"})
}

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
