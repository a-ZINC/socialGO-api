package store

import (
	"context"
	"database/sql"
	"social-api/model"
)

type PostRepo interface {
	Create(ctx context.Context, post model.Post) error
}

type UserRepo interface {
	Create(ctx context.Context, user model.User) error
}

type Store struct {
	Posts PostRepo
	Users UserRepo
}

func NewStorage(db *sql.DB) *Store {
	return &Store{
		Posts: &PostStore{
			db: db,
		},
		Users: &UserStore{
			db: db,
		},
	}
}
