package user

import (
	"database/sql"
	"errors"
)

type Service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (service *Service) CreateOrAuthUser(username string, password string) error {
	//fmt.Println("here")

	_, err := service.repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Нет пользователя - сооздадим нового
			service.repo.Save(username, password)

		}
		return err
	}

	// Пользователь есть - проверим, валиден ли пароль

	return nil
}
