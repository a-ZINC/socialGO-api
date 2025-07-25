package store

import (
	"context"
	"database/sql"
	"social-api/model"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO "User" (name, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}
