package user

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username  string
	Password  string
	Coins     int
	IsActive  bool
	CreatedAt time.Time
}
