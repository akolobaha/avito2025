package token

import "time"

type Token struct {
	UserId    int32
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Revoked   bool
}
