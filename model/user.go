package model

import (

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Password struct {
	Text *string `json:"text"`
	Hash []byte  `json:"hash"`
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Text = &password
	p.Hash = hash
	return nil
}
