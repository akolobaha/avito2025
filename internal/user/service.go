package user

import (
	"avito2015/internal/db"
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err // Обработка ошибки, если не связано с отсутствием пользователя
	}

	salt := os.Getenv("SALT")
	passwordHash := hasher.HashString(password + salt)

	if errors.Is(err, sql.ErrNoRows) {
		return uService.createUser(username, passwordHash)
	}

	return uService.authUser(user, passwordHash)
}

func (uService *Service) createUser(username string, passwordHash string) (*token.Token, error) {
	userId, err := uService.repo.Save(username, passwordHash)
	if err != nil {
		return nil, err
	}

	tRepo := token.NewTokenRepository(db.DB)
	tService := token.NewService(tRepo)

	return tService.SaveToken(userId)
}

func (uService *Service) authUser(user User, passwordHash string) (*token.Token, error) {
	if user.Password != passwordHash {
		return nil, InvalidPasswordError
	}

	tRepo := token.NewTokenRepository(db.DB)
	tService := token.NewService(tRepo)

	tokenModel, err := tService.GetByUserId(user.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return tService.SaveToken(int64(user.Id))
	}

	_, err = tService.ValidateToken(tokenModel.Jwt)
	if err != nil {
		tokenModel.IsActive = false
		if _, err := tRepo.Update(*tokenModel); err != nil {
			return nil, err
		}
		return tService.SaveToken(int64(user.Id))
	}

	return tokenModel, nil
}

func (uService *Service) GetByUsername(username string) (*User, error) {
	user, err := uService.repo.FindByUsername(username)
	if err != nil {
	}
	return &user, nil
}

func (uService *Service) GetUserAndTokenByJwt(tokenStr string) (*token.Token, *User, error) {
	repo := NewUserRepository(db.DB)

	tkn, usr, err := repo.GetByToken(tokenStr)
	if err != nil {
		return nil, nil, err
	}
	return tkn, usr, nil
}
