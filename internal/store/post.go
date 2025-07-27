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

func (s *PostStore) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO "Post" (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (model.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at, version FROM "Post"
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	post := model.Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version)
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
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
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

func (s *PostStore) Update(ctx context.Context, post model.Post, id int64) error {
	query := `
		UPDATE "Post"
		SET title=$1, content=$2, version = version + 1, updated_at = NOW()
		WHERE id = $3 And version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	p, err := s.db.ExecContext(ctx, query, post.Title, post.Content, id, post.Version)
	if err != nil {
		return err
	}
	rowsAffected, err := p.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPostNotFound
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, pagination *Pagination) ([]*model.PostWithMetadata, error) {
	query := `
		SELECT p.id, p.title, p.content, p.user_id, p.tags, p.created_at, p.updated_at, COUNT(c.id) AS comment_count
		FROM "Post" p
		LEFT JOIN "Comment" c ON c.post_id = p.id
		LEFT JOIN "User" u ON u.id = p.user_id
		LEFT JOIN "Follower" f ON f.follower_id = p.user_id AND f.user_id = $1
		WHERE (p.user_id = $1 OR f.user_id = $1) AND
		(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
		($5::timestamp IS NULL OR p.created_at >= $5::timestamp) AND
		(array_length($6::text[], 1) IS NULL OR 
       	p.tags::text[] @> $6::text[])
		GROUP BY p.id
		ORDER BY p.created_at ` + pagination.Order + `
		LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID, pagination.Limit, pagination.Offset, pagination.Search, pagination.Since, pq.Array(pagination.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*model.PostWithMetadata
	for rows.Next() {
		var post model.PostWithMetadata
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.CommentCount); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStore) GetUserFeedCount(ctx context.Context, userID int64, pagination *Pagination) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT p.id)
		FROM "Post" p
		LEFT JOIN "Comment" c ON c.post_id = p.id
		LEFT JOIN "User" u ON u.id = p.user_id
		LEFT JOIN "Follower" f ON f.follower_id = p.user_id AND f.user_id = $1
		WHERE (p.user_id = $1 OR f.user_id = $1) AND
		(p.title ILIKE '%' || $2 || '%' OR p.content ILIKE '%' || $2 || '%') AND
		($3::timestamp IS NULL OR p.created_at >= $3::timestamp) AND
		(array_length($4::text[], 1) IS NULL OR 
       	p.tags::text[] && $4::text[])
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, userID, pagination.Search, pagination.Since, pq.Array(pagination.Tags))
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
