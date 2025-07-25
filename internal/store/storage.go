package store

import (
	"context"
	"database/sql"
	"errors"
	"social-api/model"
	"time"
)

var (
	ErrPostNotFound = errors.New("resource not found")
	TimeOut         = 5 * time.Second
)

type PostRepo interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id int64) (model.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post model.Post, id int64) error
}

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
}

type CommentRepo interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByPostId(ctx context.Context, postId int64) ([]model.Comment, error)
}

type Store struct {
	Posts    PostRepo
	Users    UserRepo
	Comments CommentRepo
}

func NewStorage(db *sql.DB) *Store {
	return &Store{
		Posts: &PostStore{
			db: db,
		},
		Users: &UserStore{
			db: db,
		},
		Comments: &CommentStore{
			db: db,
		},
	}
}
