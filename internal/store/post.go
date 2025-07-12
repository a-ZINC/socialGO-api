package store

import (
	"context"
	"database/sql"
	"social-api/model"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post model.Post) error {
	query := `
		INSERT INTO "Post" (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (model.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at FROM "Post"
		WHERE id = $1
	`
	post := model.Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return post, ErrPostNotFound
		default:
			return post, err
		}
	}
	return post, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM "Post"
		WHERE id = $1
	`
	post, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := post.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrPostNotFound
    }
	return nil
}
