package token

import (
	"avito2015/internal/db"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Save(token Token) error
	Get(userID int) (*Token, error)
	Update(token Token) (*Token, error)
}

type tokenRepositoryImpl struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) Repository {
	return &tokenRepositoryImpl{db: db}
}

func (r *tokenRepositoryImpl) Save(token Token) error {
	query := `INSERT INTO user_token (jwt, user_id) VALUES ($1, $2)`
	_, err := db.DB.Exec(query, token.Jwt, token.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepositoryImpl) Get(userID int) (*Token, error) {
	var token Token
	query := `SELECT jwt, user_id, is_active FROM "user_token" WHERE user_id = $1 AND is_active = true`
	err := r.db.Get(&token, query, userID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepositoryImpl) Update(token Token) (*Token, error) {
	query := `UPDATE "user_token" SET jwt = $1, is_active = $2 WHERE user_id = $3`
	_, err := db.DB.Exec(query, token.Jwt, token.IsActive, token.UserId)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
