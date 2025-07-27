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

func (s *UserStore) GetByID(ctx context.Context, id int64) (model.User, error) {
	query := `SELECT id, name, email, created_at FROM "User" WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	var user model.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrPostNotFound
		}
		return model.User{}, err
	}
	return user, nil
}
