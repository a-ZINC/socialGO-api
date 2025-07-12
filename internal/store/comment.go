package store

import (
	"context"
	"database/sql"
	"log"
	"social-api/model"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment model.Comment) error {
	query := `
		INSERT INTO "Comment" ( content, post_id, user_id )
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := s.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]model.Comment, error) {
	query := `
		SELECT id, content, post_id, user_id, created_at FROM "Comment" c
		LEFT JOIN "User" u ON c.user_id = u.id
		WHERE c.post_id = $1
	`

	row, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var comments []model.Comment
	for row.Next() {
		var comment model.Comment
		err := row.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

