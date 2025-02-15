package merch

import (
	"avito2015/internal/user"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Buy(user *user.User, merchName string) error {
	// Получим по name модель из БД
	m, err := s.repo.FindMerchByName(merchName)
	if err != nil {
		return err
	}

	// Проверим, хватает ли денег у пользователя
	if user.Coins < m.Price {
		return errors.New("not enough coins")
	}

	// Выполним транзакцию
	err = s.repo.Buy(user, m)
	if err != nil {
		return err
	}

	return nil
}
