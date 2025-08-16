package store

import (
	"context"
	"database/sql"
	"errors"
	"social-api/model"
	"time"
)

type UserStore struct {
	db *sql.DB
}

var (
	ErrEmailTaken = errors.New("email is already taken")
	ErrNameTaken  = errors.New("name is already taken")
)

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO "User" (name, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == "unique constraint \"user_email_key\"":
			return ErrEmailTaken
		case err.Error() == "unique constraint \"user_name_key\"":
			return ErrNameTaken
		default:
			return err
		}
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

func (s *UserStore) CreateAndInvitation(ctx context.Context, user *model.User, token string, expiryTime time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		err := s.Create(ctx, tx, user)
		if err != nil {
			return err
		}
		return s.CreateInvite(ctx, tx, user, token, expiryTime)
	})
}

func (s *UserStore) CreateInvite(ctx context.Context, tx *sql.Tx, user *model.User, token string, expiryTime time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry_time)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, user.ID, time.Now().Add(expiryTime))
	if err != nil {
		return err
	}
	return nil
}
