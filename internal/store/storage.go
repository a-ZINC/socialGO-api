package store

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
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
	GetUserFeed(ctx context.Context, userID int64, pagination *Pagination) ([]*model.PostWithMetadata, error)
	GetUserFeedCount(ctx context.Context, userID int64, pagination *Pagination) (int64, error)
}

type UserRepo interface {
	Create(ctx context.Context, tx *sql.Tx, user *model.User) error
	GetByID(ctx context.Context, id int64) (model.User, error)
	CreateAndInvitation(ctx context.Context, user *model.User, token string, expiryTime time.Duration) error
	ActivateUser(ctx context.Context, token string) error
}

type CommentRepo interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByPostId(ctx context.Context, postId int64) ([]model.Comment, error)
}

type FollowerRepo interface {
	Follow(ctx context.Context, followerId, followedId int64) error
	Unfollow(ctx context.Context, followerId, followedId int64) error
}

type PaginationRepo interface {
	GetPaginated(r *http.Request) (*Pagination, error)
}
type Store struct {
	Posts      PostRepo
	Users      UserRepo
	Comments   CommentRepo
	Follower   FollowerRepo
	Pagination PaginationRepo
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
		Follower: &FollowerStore{
			db: db,
		},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
