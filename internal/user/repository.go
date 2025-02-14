package user

import (
	"avito2015/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository interface {
	Save(user string, password string) (int64, error)
	FindByUsername(username string) (User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository() Repository {
	return &userRepositoryImpl{db: db.DB}
}

func (r *userRepositoryImpl) Save(user string, password string) (int64, error) {
	query := `INSERT INTO "user" (username, password) VALUES ($1, $2) RETURNING id`
	var lastInsertID int64
	err := r.db.QueryRow(query, user, password).Scan(&lastInsertID)

	if err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (User, error) {
	var user User
	query := `SELECT id, username, coins, password, is_active, created_at FROM "user" WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		log.Println("Error querying user by username:", err)
		return User{}, err
	}

	return user, nil
}
