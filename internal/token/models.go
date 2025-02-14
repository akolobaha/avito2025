package token

import "github.com/golang-jwt/jwt/v4"

type Token struct {
	Jwt      string `json:"token" db:"jwt"`
	UserId   int64  `json:"-" db:"user_id"`
	IsActive bool   `json:"-" db:"is_active"`
}

type Claims struct {
	jwt.RegisteredClaims
}
