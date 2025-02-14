package user

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Coins     int       `db:"coins"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}
