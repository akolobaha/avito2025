package user

import (
	"avito2015/internal/token"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository interface {
	Save(user string, password string) (int64, error)
	FindByUsername(username string) (User, error)
	GetByToken(token string) (*token.Token, *User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) Repository {
	return &userRepositoryImpl{db: db}
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
	query := `SELECT id, username, coins, password, created_at FROM "user" WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		log.Println("Error querying user by username:", err)
		return User{}, err
	}

	return user, nil
}

func (r *userRepositoryImpl) GetByToken(tkn string) (*token.Token, *User, error) {
	var tokenM token.Token
	var usr User
	query := `
  SELECT ut.jwt, ut.user_id, ut.is_active, 
         u.id, u.username, u.password, u.coins, u.created_at 
  FROM user_token ut 
  JOIN public."user" u ON u.id = ut.user_id 
  WHERE ut.jwt = $1 AND ut.is_active = true
 `
	err := r.db.QueryRow(query, tkn).Scan(
		&tokenM.Jwt,
		&tokenM.UserId,
		&tokenM.IsActive,
		&usr.Id,
		&usr.Username,
		&usr.Password,
		&usr.Coins,
		&usr.CreatedAt,
	)
	if err != nil {
		return nil, nil, err
	}
	return &tokenM, &usr, nil
}
