package store

import (
	"context"
	"database/sql"
	"fmt"
)

type FollowerStore struct {
	db *sql.DB
}

func (fs *FollowerStore) Follow(ctx context.Context, userId, followedId int64) error {
	query := "INSERT INTO \"Follower\" (user_id, follower_id) VALUES ($1, $2)"
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	fmt.Printf("Following user %d with follower %d\n", userId, followedId)
	_, err := fs.db.ExecContext(ctx, query, userId, followedId)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FollowerStore) Unfollow(ctx context.Context, userId, followedId int64) error {
	query := "DELETE FROM \"Follower\" WHERE user_id = $1 AND follower_id = $2"
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	_, err := fs.db.ExecContext(ctx, query, userId, followedId)
	if err != nil {
		return err
	}
	return nil
}
