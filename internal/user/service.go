package user

import (
	"avito2015/internal/token"
	"avito2015/pkg/hasher"
	"database/sql"
	"errors"
	"os"
)

type Service struct {
	repo Repository
}

var InvalidPasswordError = errors.New("invalid password")

func NewUserService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (uService *Service) CreateOrAuthUser(username string, password string) (*token.Token, error) {
	user, err := uService.repo.FindByUsername(username)
	salt := os.Getenv("SALT")
	passwordHash := hasher.HashString(password + salt)

	// Создаем нового
	if errors.Is(err, sql.ErrNoRows) {
		userId, err := uService.repo.Save(username, passwordHash)

		if err != nil {
			return nil, err
		}

		token, err := token.SaveToken(userId)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	// Авторизуем существующего
	if user.Password == passwordHash {
		tokenRepo := token.NewTokenRepository()
		tokenModel, err := token.GetByUserId(user.Id)

		_, err = token.ValidateToken(tokenModel.Jwt)

		if err != nil {
			// Инвалидируем старые токен
			tokenModel.IsActive = false
			_, err := tokenRepo.Update(*tokenModel)
			if err != nil {
				return nil, err
			}

			// Создаем новый
			tkn, err := token.SaveToken(int64(user.Id))
			return tkn, nil
		}

		return tokenModel, nil
	}

	return nil, InvalidPasswordError
}

func (uService *Service) GetByUsername(username string) (*User, error) {
	user, err := uService.repo.FindByUsername(username)
	if err != nil {
	}
	return &user, nil
}

func (uService *Service) GetUserAndTokenByJwt(tokenStr string) (*token.Token, *User, error) {
	repo := NewUserRepository()

	tkn, usr, err := repo.GetByToken(tokenStr)
	if err != nil {
		return nil, nil, err
	}
	return tkn, usr, nil
}
