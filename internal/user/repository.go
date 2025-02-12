package user

import (
	"avito2015/internal/db"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type UserRepository interface {
	Save(user string, password string) error
	FindByUsername(username string) (User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{db: db.DB}
}

func (r *userRepositoryImpl) Save(user string, password string) error {
	// Добавить логику хэширования пароля

	query := `INSERT INTO "user" (username, password) VALUES ($1, $2)`
	_, err := r.db.Exec(query, user, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (User, error) {
	var user User
	query := `SELECT username, coins FROM "user" WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.Username, &user.Coins)
	if err != nil {
		log.Println("Error querying user by username:", err)
		return User{}, err
	}

	return user, nil
}
