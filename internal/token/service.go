package token

import (
	"avito2015/internal/db"
	jwt "github.com/golang-jwt/jwt/v4"
	"os"
	"strconv"
	"time"
)

var tokenExpPeriod time.Duration = 24 * time.Hour

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func generateToken(userId int64) (string, error) {
	expirationTime := time.Now().Add(tokenExpPeriod)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.FormatInt(userId, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *Service) SaveToken(userId int64) (*Token, error) {
	var token Token

	jwtStr, err := generateToken(userId)
	if err != nil {
		return nil, err
	}
	token.Jwt = jwtStr
	token.UserId = userId
	token.IsActive = true

	err = s.repo.Save(token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *Service) GetByUserId(userId int) (*Token, error) {
	repo := NewTokenRepository(db.DB)
	token, err := repo.Get(userId)
	if err != nil {
		return nil, err
	}

	return token, nil
}
