package store

import (
	"context"
	"database/sql"
	"social-api/model"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post model.Post) error {
	query := `
		INSERT INTO posts (title, content, userId)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err;
	}
	return nil
}
